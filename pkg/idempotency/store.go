package idempotency

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrDuplicateRequest = errors.New("duplicate request")

type Store struct {
	rdb *redis.Client
}

func NewStore(rdb *redis.Client) *Store {
	return &Store{rdb: rdb}
}

func (s *Store) Reserve(ctx context.Context, key string, ttl time.Duration) error {
	ok, err := s.rdb.SetNX(ctx, key, "1", ttl).Result()
	if err != nil {
		return err
	}
	if !ok {
		return ErrDuplicateRequest
	}
	return nil
}
