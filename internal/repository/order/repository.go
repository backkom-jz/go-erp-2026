package order

import (
	"context"
	"go-erp/internal/domain/order"

	"gorm.io/gorm"
)

type Repository interface {
	// Create 创建订单主表及订单项（事务内调用）。
	Create(ctx context.Context, tx *gorm.DB, entity *order.Order, items []order.OrderItem) error
	// GetByID 按订单 ID 查询订单及订单项。
	GetByID(ctx context.Context, id uint) (*order.Order, []order.OrderItem, error)
	// UpdateStatusByOrderNo 按订单号更新状态。
	UpdateStatusByOrderNo(ctx context.Context, orderNo string, status string) error
	// CancelIfPending 仅当订单为 pending 时取消。
	CancelIfPending(ctx context.Context, orderNo string) (bool, error)
	// MarkPaidPreferPaid 将订单置为 paid（支持 cancelled 回正）。
	MarkPaidPreferPaid(ctx context.Context, orderNo string) (bool, error)
}

type GormRepository struct {
	db *gorm.DB
}

// NewRepository 创建订单仓储实现。
func NewRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

// Create 创建订单及订单项。
func (r *GormRepository) Create(ctx context.Context, tx *gorm.DB, entity *order.Order, items []order.OrderItem) error {
	if err := tx.WithContext(ctx).Create(entity).Error; err != nil {
		return err
	}
	for i := range items {
		items[i].OrderID = entity.ID
	}
	return tx.WithContext(ctx).Create(&items).Error
}

// GetByID 查询订单详情。
func (r *GormRepository) GetByID(ctx context.Context, id uint) (*order.Order, []order.OrderItem, error) {
	var header order.Order
	if err := r.db.WithContext(ctx).First(&header, id).Error; err != nil {
		return nil, nil, err
	}
	var items []order.OrderItem
	if err := r.db.WithContext(ctx).Where("order_id = ?", id).Find(&items).Error; err != nil {
		return nil, nil, err
	}
	return &header, items, nil
}

// UpdateStatusByOrderNo 按订单号更新状态。
func (r *GormRepository) UpdateStatusByOrderNo(ctx context.Context, orderNo string, status string) error {
	return r.db.WithContext(ctx).Model(&order.Order{}).
		Where("order_no = ?", orderNo).
		Update("status", status).Error
}

// CancelIfPending 仅取消 pending 状态订单。
func (r *GormRepository) CancelIfPending(ctx context.Context, orderNo string) (bool, error) {
	result := r.db.WithContext(ctx).Model(&order.Order{}).
		Where("order_no = ? AND status = ?", orderNo, order.StatusPending).
		Update("status", order.StatusCancelled)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

// MarkPaidPreferPaid 将 pending/cancelled 订单置为 paid。
func (r *GormRepository) MarkPaidPreferPaid(ctx context.Context, orderNo string) (bool, error) {
	result := r.db.WithContext(ctx).
		Model(&order.Order{}).
		Where("order_no = ? AND status IN ?", orderNo, []string{order.StatusPending, order.StatusCancelled}).
		Update("status", order.StatusPaid)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
