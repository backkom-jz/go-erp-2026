package payment

import (
	"context"
	"errors"
	"go-erp/internal/domain/payment"

	"gorm.io/gorm"
)

type Repository interface {
	CreateOrUpdateByPaymentNo(ctx context.Context, req payment.Record) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) CreateOrUpdateByPaymentNo(ctx context.Context, req payment.Record) error {
	var existing payment.Record
	err := r.db.WithContext(ctx).Where("payment_no = ?", req.PaymentNo).First(&existing).Error
	if err == nil {
		existing.Status = req.Status
		existing.CallbackRaw = req.CallbackRaw
		existing.Channel = req.Channel
		return r.db.WithContext(ctx).Save(&existing).Error
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return r.db.WithContext(ctx).Create(&req).Error
}
