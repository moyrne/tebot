package ratelimit

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
	"github.com/pkg/errors"
)

type Limiter struct {
	limiter *redis_rate.Limiter
}

func NewRedisLimit(client *redis.Client) *Limiter {
	return &Limiter{limiter: redis_rate.NewLimiter(client)}
}

var ErrRateLimit = errors.New("rate remaining is zero")

// Rate TODO 限流 防止封号 (2t/min CD)
func (l *Limiter) Rate(ctx context.Context, name string, userID string) error {
	// auto_reply_rate_limit_
	result, err := l.limiter.Allow(ctx, name+"_"+userID, redis_rate.PerMinute(3))
	if err != nil {
		return errors.WithStack(err)
	}
	if result.Remaining == 0 {
		return errors.WithStack(ErrRateLimit)
	}
	return nil
}
