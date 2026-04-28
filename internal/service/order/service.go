package order

import (
	"context"
	"encoding/json"
	domainorder "go-erp/internal/domain/order"
	dtoinventory "go-erp/internal/dto/inventory"
	dtoorder "go-erp/internal/dto/order"
	inventorysvc "go-erp/internal/service/inventory"
	orderrepo "go-erp/internal/repository/order"
	"go-erp/pkg/dbtx"
	"go-erp/pkg/errs"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TimeoutPayload struct {
	OrderNo string `json:"order_no"`
}

type Service struct {
	db        *gorm.DB
	repo      orderrepo.Repository
	outbox    orderrepo.OutboxRepository
	inventory *inventorysvc.Service
}

// NewService 创建订单服务。
func NewService(
	db *gorm.DB,
	repo orderrepo.Repository,
	outbox orderrepo.OutboxRepository,
	inventory *inventorysvc.Service,
) *Service {
	return &Service{
		db:        db,
		repo:      repo,
		outbox:    outbox,
		inventory: inventory,
	}
}

// Create 创建订单。
// 备注：订单创建与库存扣减、Outbox 写入在同一事务内完成。
func (s *Service) Create(ctx context.Context, req dtoorder.CreateOrderRequest) (*domainorder.Order, error) {
	orderNo := uuid.NewString()
	items := make([]domainorder.OrderItem, 0, len(req.Items))
	var total int64

	for _, it := range req.Items {
		items = append(items, domainorder.OrderItem{
			SKUID:      it.SKUID,
			Qty:        it.Qty,
			PriceCents: it.PriceCents,
		})
		total += it.PriceCents * it.Qty
	}

	header := &domainorder.Order{
		OrderNo:    orderNo,
		UserID:     req.UserID,
		TenantID:   req.TenantID,
		Status:     domainorder.StatusPending,
		TotalCents: total,
	}

	err := dbtx.WithTransaction(ctx, s.db, func(tx *gorm.DB) error {
		for _, item := range req.Items {
			if err := s.inventory.DeductWithTx(ctx, tx, dtoinventory.DeductRequest{
				SKUID:      item.SKUID,
				Qty:        item.Qty,
				BusinessNo: orderNo,
			}); err != nil {
				return err
			}
		}
		if err := s.repo.Create(ctx, tx, header, items); err != nil {
			return err
		}

		createdPayload, _ := json.Marshal(map[string]interface{}{
			"order_no": orderNo,
			"type":     "order.created",
		})
		if err := s.outbox.CreateInTx(ctx, tx, &domainorder.OutboxEvent{
			EventType:  "order.created",
			RoutingKey: "order.created",
			Payload:    string(createdPayload),
			Status:     domainorder.OutboxStatusPending,
		}); err != nil {
			return err
		}

		timeoutPayload, _ := json.Marshal(TimeoutPayload{OrderNo: orderNo})
		if err := s.outbox.CreateInTx(ctx, tx, &domainorder.OutboxEvent{
			EventType:  "order.timeout.delay",
			RoutingKey: "order.timeout.delay",
			Payload:    string(timeoutPayload),
			Status:     domainorder.OutboxStatusPending,
		}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return header, nil
}

// GetByID 查询订单详情。
func (s *Service) GetByID(ctx context.Context, id uint) (*domainorder.Order, []domainorder.OrderItem, error) {
	header, items, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, errs.Wrap(errs.CodeNotFound, "order_not_found", err)
	}
	return header, items, nil
}

// MarkPaid 将订单标记为已支付。
func (s *Service) MarkPaid(ctx context.Context, orderNo string) error {
	_, err := s.repo.MarkPaidPreferPaid(ctx, orderNo)
	return err
}

// CancelIfTimeout 超时取消订单（仅 pending 状态可取消）。
func (s *Service) CancelIfTimeout(ctx context.Context, orderNo string) error {
	_, err := s.repo.CancelIfPending(ctx, orderNo)
	return err
}

// HandleTimeoutMessage 处理订单超时消息。
func (s *Service) HandleTimeoutMessage(ctx context.Context, payload []byte) error {
	var msg TimeoutPayload
	if err := json.Unmarshal(payload, &msg); err != nil {
		return errs.Wrap(errs.CodeBadRequest, "invalid_timeout_payload", err)
	}
	if msg.OrderNo == "" {
		return errs.New(errs.CodeBadRequest, "missing_order_no")
	}
	return s.CancelIfTimeout(ctx, msg.OrderNo)
}

type OutboxDispatcher struct {
	repo      orderrepo.OutboxRepository
	publisher interface {
		Publish(context.Context, string, []byte) error
	}
	interval   time.Duration
	maxRetry   int
	baseBackoff time.Duration
}

// NewOutboxDispatcher 创建 Outbox 派发器。
func NewOutboxDispatcher(
	repo orderrepo.OutboxRepository,
	publisher interface {
		Publish(context.Context, string, []byte) error
	},
	interval time.Duration,
	maxRetry int,
	baseBackoff time.Duration,
) *OutboxDispatcher {
	if interval <= 0 {
		interval = 2 * time.Second
	}
	if maxRetry <= 0 {
		maxRetry = 5
	}
	if baseBackoff <= 0 {
		baseBackoff = 3 * time.Second
	}
	return &OutboxDispatcher{
		repo:        repo,
		publisher:   publisher,
		interval:    interval,
		maxRetry:    maxRetry,
		baseBackoff: baseBackoff,
	}
}

// Run 启动 Outbox 派发轮询。
func (d *OutboxDispatcher) Run(ctx context.Context) {
	if d.publisher == nil || d.repo == nil {
		return
	}
	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()

	d.dispatchOnce(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			d.dispatchOnce(ctx)
		}
	}
}

func (d *OutboxDispatcher) dispatchOnce(ctx context.Context) {
	events, err := d.repo.FetchPending(ctx, 100)
	if err != nil {
		return
	}
	for _, evt := range events {
		if err := d.publisher.Publish(ctx, evt.RoutingKey, []byte(evt.Payload)); err != nil {
			if evt.RetryCount+1 >= d.maxRetry {
				_ = d.repo.MarkDead(ctx, evt, err.Error())
				continue
			}
			retryDelay := d.baseBackoff * time.Duration(1<<evt.RetryCount)
			_ = d.repo.MarkRetry(ctx, evt.ID, time.Now().Add(retryDelay))
			continue
		}
		_ = d.repo.MarkSent(ctx, evt.ID)
	}
}
