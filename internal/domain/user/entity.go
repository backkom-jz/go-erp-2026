package user

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserNo    string    `gorm:"size:64;uniqueIndex" json:"user_no"`
	Name      string    `gorm:"size:128" json:"name"`
	TenantID  string    `gorm:"size:64;index" json:"tenant_id"`
	Role      string    `gorm:"size:32" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
