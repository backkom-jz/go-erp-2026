package order

import "time"

type OutboxDeadLetter struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	OutboxID   uint      `gorm:"index" json:"outbox_id"`
	EventType  string    `gorm:"size:128;index" json:"event_type"`
	RoutingKey string    `gorm:"size:128;index" json:"routing_key"`
	Payload    string    `gorm:"type:text" json:"payload"`
	ErrorMsg   string    `gorm:"size:1024" json:"error_msg"`
	CreatedAt  time.Time `json:"created_at"`
}
