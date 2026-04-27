package bootstrap

import (
	"go-erp/pkg/mq"
	"go-erp/pkg/mq/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	orderTimeoutDelayQueue   = "order.timeout.delay.q"
	orderTimeoutProcessQueue = "order.timeout.process.q"
)

type MQClient struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Publisher  mq.Publisher
}

func InitMQ(cfg MQConfig) (*MQClient, error) {
	if !cfg.Enabled {
		return nil, nil
	}
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	if err := ch.ExchangeDeclare(cfg.Exchange, "topic", true, false, false, false, nil); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, err
	}
	if err := declareOrderTimeoutTopology(ch, cfg); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, err
	}
	return &MQClient{
		Connection: conn,
		Channel:    ch,
		Publisher:  rabbitmq.NewPublisher(ch, cfg.Exchange),
	}, nil
}

func declareOrderTimeoutTopology(ch *amqp.Channel, cfg MQConfig) error {
	delayArgs := amqp.Table{
		"x-message-ttl":             int32(cfg.OrderTimeoutMinutes * 60 * 1000),
		"x-dead-letter-exchange":    cfg.Exchange,
		"x-dead-letter-routing-key": "order.timeout.process",
	}
	if _, err := ch.QueueDeclare(orderTimeoutDelayQueue, true, false, false, false, delayArgs); err != nil {
		return err
	}
	if err := ch.QueueBind(orderTimeoutDelayQueue, "order.timeout.delay", cfg.Exchange, false, nil); err != nil {
		return err
	}

	if _, err := ch.QueueDeclare(orderTimeoutProcessQueue, true, false, false, false, nil); err != nil {
		return err
	}
	if err := ch.QueueBind(orderTimeoutProcessQueue, "order.timeout.process", cfg.Exchange, false, nil); err != nil {
		return err
	}
	return nil
}
