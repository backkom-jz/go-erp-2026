package order

import "time"

const (
	StatusPending   = "pending_payment"
	StatusPaid      = "paid"
	StatusShipped   = "shipped"
	StatusCompleted = "completed"
	StatusCancelled = "cancelled"
)

type Order struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	OrderNo    string    `gorm:"size:64;uniqueIndex" json:"order_no"`
	UserID     uint      `gorm:"index" json:"user_id"`
	TenantID   string    `gorm:"size:64;index" json:"tenant_id"`
	Status     string    `gorm:"size:32;index" json:"status"`
	TotalCents int64     `json:"total_cents"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type OrderItem struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	OrderID    uint      `gorm:"index" json:"order_id"`
	SKUID      uint      `json:"sku_id"`
	Qty        int64     `json:"qty"`
	PriceCents int64     `json:"price_cents"`
	CreatedAt  time.Time `json:"created_at"`
}
