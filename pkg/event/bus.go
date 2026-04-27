package event

import (
	"context"
	"sync"
)

type Handler func(ctx context.Context, payload []byte) error

type Bus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

func NewBus() *Bus {
	return &Bus{handlers: make(map[string][]Handler)}
}

func (b *Bus) Subscribe(topic string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[topic] = append(b.handlers[topic], handler)
}

func (b *Bus) Publish(ctx context.Context, topic string, payload []byte) error {
	b.mu.RLock()
	handlers := b.handlers[topic]
	b.mu.RUnlock()
	for _, h := range handlers {
		if err := h(ctx, payload); err != nil {
			return err
		}
	}
	return nil
}
