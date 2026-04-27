package mq

import "context"

type Publisher interface {
	Publish(ctx context.Context, routingKey string, payload []byte) error
}
