package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ch       *amqp.Channel
	exchange string
}

func NewPublisher(ch *amqp.Channel, exchange string) *Publisher {
	return &Publisher{ch: ch, exchange: exchange}
}

func (p *Publisher) Publish(ctx context.Context, routingKey string, payload []byte) error {
	return p.ch.PublishWithContext(ctx, p.exchange, routingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        payload,
	})
}
