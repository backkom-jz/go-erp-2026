package product

import "time"

type SPU struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"size:256;index" json:"name"`
	CategoryID uint      `json:"category_id"`
	Brand      string    `gorm:"size:128" json:"brand"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type SKU struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	SPUID      uint      `gorm:"index" json:"spu_id"`
	Code       string    `gorm:"size:64;uniqueIndex" json:"code"`
	Name       string    `gorm:"size:256" json:"name"`
	PriceCents int64     `json:"price_cents"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:128;uniqueIndex" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
