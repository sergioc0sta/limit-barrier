package storage

import (
	"context"

	"github.com/sergioc0sta/limit-barrier/internal/limiter"
)

type Store interface {
	Ping(ctx context.Context) error
	Set(ctx context.Context, key, value string) error
	Get(ctx context.Context, key string) (string, error)
	Close() error
	limiter.RateLimiterStrategy
}
