package product

import (
	"context"
	"go-erp/internal/domain/product"

	"gorm.io/gorm"
)

type Repository interface {
	CreateSPU(ctx context.Context, entity *product.SPU) error
	ListSPU(ctx context.Context, limit int) ([]product.SPU, error)
	CreateSKU(ctx context.Context, entity *product.SKU) error
	GetSKU(ctx context.Context, id uint) (*product.SKU, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) CreateSPU(ctx context.Context, entity *product.SPU) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *GormRepository) ListSPU(ctx context.Context, limit int) ([]product.SPU, error) {
	var rows []product.SPU
	if limit <= 0 {
		limit = 20
	}
	if err := r.db.WithContext(ctx).Order("id desc").Limit(limit).Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *GormRepository) CreateSKU(ctx context.Context, entity *product.SKU) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *GormRepository) GetSKU(ctx context.Context, id uint) (*product.SKU, error) {
	var row product.SKU
	if err := r.db.WithContext(ctx).First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}
