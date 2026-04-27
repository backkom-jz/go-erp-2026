package payment

import "time"

type Record struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	OrderNo     string    `gorm:"size:64;index" json:"order_no"`
	PaymentNo   string    `gorm:"size:64;uniqueIndex" json:"payment_no"`
	Channel     string    `gorm:"size:32" json:"channel"`
	Status      string    `gorm:"size:32;index" json:"status"`
	CallbackRaw string    `gorm:"type:text" json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
