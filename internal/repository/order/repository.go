package order

import (
	"context"
	"go-erp/internal/domain/order"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, tx *gorm.DB, entity *order.Order, items []order.OrderItem) error
	GetByID(ctx context.Context, id uint) (*order.Order, []order.OrderItem, error)
	UpdateStatusByOrderNo(ctx context.Context, orderNo string, status string) error
	CancelIfPending(ctx context.Context, orderNo string) (bool, error)
	MarkPaidPreferPaid(ctx context.Context, orderNo string) (bool, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(ctx context.Context, tx *gorm.DB, entity *order.Order, items []order.OrderItem) error {
	if err := tx.WithContext(ctx).Create(entity).Error; err != nil {
		return err
	}
	for i := range items {
		items[i].OrderID = entity.ID
	}
	return tx.WithContext(ctx).Create(&items).Error
}

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

func (r *GormRepository) UpdateStatusByOrderNo(ctx context.Context, orderNo string, status string) error {
	return r.db.WithContext(ctx).Model(&order.Order{}).
		Where("order_no = ?", orderNo).
		Update("status", status).Error
}

func (r *GormRepository) CancelIfPending(ctx context.Context, orderNo string) (bool, error) {
	result := r.db.WithContext(ctx).Model(&order.Order{}).
		Where("order_no = ? AND status = ?", orderNo, order.StatusPending).
		Update("status", order.StatusCancelled)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

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
