package payment

import (
	"context"
	"encoding/json"
	domainpayment "go-erp/internal/domain/payment"
	dtopayment "go-erp/internal/dto/payment"
	orderService "go-erp/internal/service/order"
	paymentrepo "go-erp/internal/repository/payment"
	"go-erp/pkg/errs"
	"go-erp/pkg/idempotency"
	"time"
)

type Service struct {
	repo      paymentrepo.Repository
	orders    *orderService.Service
	idemStore *idempotency.Store
}

func NewService(repo paymentrepo.Repository, orders *orderService.Service, idemStore *idempotency.Store) *Service {
	return &Service{
		repo:      repo,
		orders:    orders,
		idemStore: idemStore,
	}
}

func (s *Service) Callback(ctx context.Context, req dtopayment.CallbackRequest) error {
	if s.idemStore != nil {
		if err := s.idemStore.Reserve(ctx, "pay:"+req.PaymentNo, 10*time.Minute); err != nil {
			return errs.New(errs.CodeDuplicate, "duplicate_payment_callback")
		}
	}
	raw, _ := json.Marshal(req)
	if err := s.repo.CreateOrUpdateByPaymentNo(ctx, domainpayment.Record{
		OrderNo:     req.OrderNo,
		PaymentNo:   req.PaymentNo,
		Channel:     req.Channel,
		Status:      req.Status,
		CallbackRaw: string(raw),
	}); err != nil {
		return errs.Wrap(errs.CodeInternal, "save_payment_failed", err)
	}

	if req.Status == "paid" {
		if err := s.orders.MarkPaid(ctx, req.OrderNo); err != nil {
			return errs.Wrap(errs.CodeInternal, "mark_order_paid_failed", err)
		}
	}
	return nil
}
