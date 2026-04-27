package limit

import (
	"errors"
	"sync"
	"time"
)

var ErrRateLimited = errors.New("rate limited")

type WindowLimiter struct {
	mu      sync.Mutex
	window  time.Duration
	maxHits int
	hits    map[string][]time.Time
}

func NewWindowLimiter(window time.Duration, maxHits int) *WindowLimiter {
	return &WindowLimiter{
		window:  window,
		maxHits: maxHits,
		hits:    make(map[string][]time.Time),
	}
}

func (l *WindowLimiter) Allow(key string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	records := l.hits[key]
	valid := records[:0]
	for _, t := range records {
		if now.Sub(t) <= l.window {
			valid = append(valid, t)
		}
	}
	if len(valid) >= l.maxHits {
		l.hits[key] = valid
		return ErrRateLimited
	}
	l.hits[key] = append(valid, now)
	return nil
}
