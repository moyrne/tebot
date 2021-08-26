package ratelimit

import (
	"context"
	"sync"
)

var (
	mu   sync.RWMutex
	rate RateLimit
)

type RateLimit interface {
	Rate(ctx context.Context, name, userID string) error
}

func InitRate(limit RateLimit) {
	mu.Lock()
	defer mu.Unlock()
	rate = limit
}

func Rate(ctx context.Context, name, userID string) error {
	mu.RLocker()
	defer mu.RUnlock()
	return rate.Rate(ctx, name, userID)
}
