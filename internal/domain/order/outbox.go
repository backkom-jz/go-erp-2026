package order

import "time"

const (
	OutboxStatusPending = "pending"
	OutboxStatusSent    = "sent"
	OutboxStatusFailed  = "failed"
	OutboxStatusDead    = "dead"
)

type OutboxEvent struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	EventType   string     `gorm:"size:128;index" json:"event_type"`
	RoutingKey  string     `gorm:"size:128;index" json:"routing_key"`
	Payload     string     `gorm:"type:text" json:"payload"`
	Status      string     `gorm:"size:32;index" json:"status"`
	RetryCount  int        `json:"retry_count"`
	NextRetryAt *time.Time `gorm:"index" json:"next_retry_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
