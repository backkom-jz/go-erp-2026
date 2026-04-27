package bootstrap

import (
	"go-erp/internal/domain/inventory"
	"go-erp/internal/domain/order"
	"go-erp/internal/domain/payment"
	"go-erp/internal/domain/product"
	"go-erp/internal/domain/user"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&user.User{},
		&product.SPU{},
		&product.SKU{},
		&product.Category{},
		&inventory.Inventory{},
		&inventory.InventoryLog{},
		&order.Order{},
		&order.OrderItem{},
		&order.OutboxEvent{},
		&order.OutboxDeadLetter{},
		&payment.Record{},
	)
}
