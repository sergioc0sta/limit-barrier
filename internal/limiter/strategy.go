package limiter

import "time"

type RateLimiterStrategy interface {
	Allow(key string, maxReq int, duration time.Duration) (allowed bool, remaining int, reset time.Time, err error)
	Block(key string, duration time.Duration) error
	Unblock(key string) error
}
