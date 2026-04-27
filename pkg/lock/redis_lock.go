package lock

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var ErrLockNotAcquired = errors.New("lock not acquired")

type RedisLocker struct {
	rdb *redis.Client
}

func NewRedisLocker(rdb *redis.Client) *RedisLocker {
	return &RedisLocker{rdb: rdb}
}

type Handle struct {
	key   string
	token string
	rdb   *redis.Client
}

func (l *RedisLocker) Acquire(ctx context.Context, key string, ttl time.Duration) (*Handle, error) {
	token := uuid.NewString()
	ok, err := l.rdb.SetNX(ctx, key, token, ttl).Result()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrLockNotAcquired
	}
	return &Handle{key: key, token: token, rdb: l.rdb}, nil
}

func (h *Handle) Release(ctx context.Context) error {
	script := redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
end
return 0
`)
	_, err := script.Run(ctx, h.rdb, []string{h.key}, h.token).Result()
	return err
}
