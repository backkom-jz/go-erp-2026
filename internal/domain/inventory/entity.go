package inventory

import "time"

type Inventory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SKUID     uint      `gorm:"uniqueIndex" json:"sku_id"`
	Stock     int64     `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type InventoryLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	SKUID      uint      `gorm:"index" json:"sku_id"`
	ChangeQty  int64     `json:"change_qty"`
	BeforeQty  int64     `json:"before_qty"`
	AfterQty   int64     `json:"after_qty"`
	BusinessNo string    `gorm:"size:64;index" json:"business_no"`
	CreatedAt  time.Time `json:"created_at"`
}
