package cqhttp

import (
	"context"
	"strconv"

	"github.com/go-redis/redis_rate/v9"
	"github.com/moyrne/tebot/internal/database"
	"github.com/pkg/errors"
)

type Limiter struct {
	limiter *redis_rate.Limiter
}

var rateLimiter *Limiter

func InitLimiter() {
	rateLimiter = &Limiter{limiter: redis_rate.NewLimiter(database.Redis)}
}

var ErrRateLimit = errors.New("rate remaining is zero")

// Rate TODO 限流 防止封号 (2t/min CD)
func (l *Limiter) Rate(ctx context.Context, name string, quid int) error {
	result, err := l.limiter.Allow(ctx, "analyze_"+name+"_"+strconv.Itoa(quid), redis_rate.PerMinute(3))
	if err != nil {
		return errors.WithStack(err)
	}
	if result.Remaining == 0 {
		return errors.WithStack(ErrRateLimit)
	}
	return nil
}
