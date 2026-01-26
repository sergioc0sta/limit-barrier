package limiter

import (
	"errors"
	"time"
)

type Limiter struct {
	strategy    RateLimiterStrategy
	ipMaxReq    int
	tokenMaxReq int
	window      time.Duration
	blockTime   time.Duration
}

func NewLimiter(strategy RateLimiterStrategy, ipMaxReq, tokenMaxReq int, window, blockTime time.Duration) (*Limiter, error) {
	if strategy == nil {
		return nil, errors.New("strategy is required")
	}
	if ipMaxReq <= 0 {
		return nil, errors.New("ipMaxReq must be > 0")
	}
	if tokenMaxReq <= 0 {
		return nil, errors.New("tokenMaxReq must be > 0")
	}
	if window <= 0 {
		return nil, errors.New("window must be > 0")
	}
	if blockTime <= 0 {
		return nil, errors.New("blockTime must be > 0")
	}

	return &Limiter{
		strategy:    strategy,
		ipMaxReq:    ipMaxReq,
		tokenMaxReq: tokenMaxReq,
		window:      window,
		blockTime:   blockTime,
	}, nil
}

func (l *Limiter) Check(key string, maxReq int) (bool, int, time.Time, error) {
	if l == nil {
		return false, 0, time.Time{}, errors.New("limiter is nil")
	}
	if key == "" {
		return false, 0, time.Time{}, errors.New("key is required")
	}
	if maxReq <= 0 {
		return false, 0, time.Time{}, errors.New("maxReq must be > 0")
	}

	allowed, remaining, reset, err := l.strategy.Allow(key, maxReq, l.window)
	if err != nil {
		return false, 0, time.Time{}, err
	}
	if !allowed {
		if blockErr := l.strategy.Block(key, l.blockTime); blockErr != nil {
			return false, 0, reset, blockErr
		}
	}

	return allowed, remaining, reset, nil
}

func (l *Limiter) IPMaxReq() int {
	return l.ipMaxReq
}

func (l *Limiter) TokenMaxReq() int {
	return l.tokenMaxReq
}

func (l *Limiter) Window() time.Duration {
	return l.window
}

func (l *Limiter) BlockTime() time.Duration {
	return l.blockTime
}
