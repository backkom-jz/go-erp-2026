package event

import (
	"context"
	"testing"
)

func TestBusPublish(t *testing.T) {
	b := NewBus()
	called := false
	b.Subscribe("topic", func(ctx context.Context, payload []byte) error {
		called = true
		return nil
	})
	if err := b.Publish(context.Background(), "topic", []byte("ok")); err != nil {
		t.Fatalf("publish failed: %v", err)
	}
	if !called {
		t.Fatalf("handler not called")
	}
}
