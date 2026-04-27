package inventory

import (
	"context"
	"go-erp/internal/domain/inventory"

	"gorm.io/gorm"
)

type Repository interface {
	GetBySKUID(ctx context.Context, skuID uint) (*inventory.Inventory, error)
	Deduct(ctx context.Context, tx *gorm.DB, skuID uint, qty int64, businessNo string) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) GetBySKUID(ctx context.Context, skuID uint) (*inventory.Inventory, error) {
	var row inventory.Inventory
	if err := r.db.WithContext(ctx).Where("sku_id = ?", skuID).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *GormRepository) Deduct(ctx context.Context, tx *gorm.DB, skuID uint, qty int64, businessNo string) error {
	var current inventory.Inventory
	if err := tx.WithContext(ctx).Where("sku_id = ?", skuID).First(&current).Error; err != nil {
		return err
	}
	if err := tx.WithContext(ctx).
		Model(&inventory.Inventory{}).
		Where("sku_id = ? AND stock >= ?", skuID, qty).
		Update("stock", gorm.Expr("stock - ?", qty)).Error; err != nil {
		return err
	}
	var latest inventory.Inventory
	if err := tx.WithContext(ctx).Where("sku_id = ?", skuID).First(&latest).Error; err != nil {
		return err
	}
	log := inventory.InventoryLog{
		SKUID:      skuID,
		ChangeQty:  -qty,
		BeforeQty:  current.Stock,
		AfterQty:   latest.Stock,
		BusinessNo: businessNo,
	}
	return tx.WithContext(ctx).Create(&log).Error
}
