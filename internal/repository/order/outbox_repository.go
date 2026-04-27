package order

import (
	"context"
	"go-erp/internal/domain/order"
	"time"

	"gorm.io/gorm"
)

type OutboxRepository interface {
	CreateInTx(ctx context.Context, tx *gorm.DB, event *order.OutboxEvent) error
	FetchPending(ctx context.Context, limit int) ([]order.OutboxEvent, error)
	MarkSent(ctx context.Context, id uint) error
	MarkRetry(ctx context.Context, id uint, retryAt time.Time) error
	MarkDead(ctx context.Context, evt order.OutboxEvent, reason string) error
}

func (r *GormRepository) CreateInTx(ctx context.Context, tx *gorm.DB, event *order.OutboxEvent) error {
	return tx.WithContext(ctx).Create(event).Error
}

func (r *GormRepository) FetchPending(ctx context.Context, limit int) ([]order.OutboxEvent, error) {
	if limit <= 0 {
		limit = 50
	}
	now := time.Now()
	var rows []order.OutboxEvent
	err := r.db.WithContext(ctx).
		Where("status = ?", order.OutboxStatusPending).
		Where("next_retry_at IS NULL OR next_retry_at <= ?", now).
		Order("id asc").
		Limit(limit).
		Find(&rows).Error
	return rows, err
}

func (r *GormRepository) MarkSent(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&order.OutboxEvent{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       order.OutboxStatusSent,
			"next_retry_at": nil,
		}).Error
}

func (r *GormRepository) MarkRetry(ctx context.Context, id uint, retryAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&order.OutboxEvent{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":        order.OutboxStatusPending,
			"retry_count":   gorm.Expr("retry_count + 1"),
			"next_retry_at": retryAt,
		}).Error
}

func (r *GormRepository) MarkDead(ctx context.Context, evt order.OutboxEvent, reason string) error {
	if len(reason) > 1024 {
		reason = reason[:1024]
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).
			Model(&order.OutboxEvent{}).
			Where("id = ?", evt.ID).
			Updates(map[string]interface{}{
				"status":        order.OutboxStatusDead,
				"next_retry_at": nil,
			}).Error; err != nil {
			return err
		}
		return tx.WithContext(ctx).Create(&order.OutboxDeadLetter{
			OutboxID:   evt.ID,
			EventType:  evt.EventType,
			RoutingKey: evt.RoutingKey,
			Payload:    evt.Payload,
			ErrorMsg:   reason,
		}).Error
	})
}
