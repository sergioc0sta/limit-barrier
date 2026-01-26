package limiter

import (
	"errors"
	"testing"
	"time"
)

type stubStrategy struct {
	allowAllowed   bool
	allowRemaining int
	allowReset     time.Time
	allowErr       error

	lastKey      string
	lastMaxReq   int
	lastDuration time.Duration

	blockCalls int
	blockErr   error
}

func (s *stubStrategy) Allow(key string, maxReq int, duration time.Duration) (bool, int, time.Time, error) {
	s.lastKey = key
	s.lastMaxReq = maxReq
	s.lastDuration = duration
	return s.allowAllowed, s.allowRemaining, s.allowReset, s.allowErr
}

func (s *stubStrategy) Block(key string, duration time.Duration) error {
	s.blockCalls++
	return s.blockErr
}

func (s *stubStrategy) Unblock(key string) error {
	return nil
}

func TestLimiterCheck_Allowed(t *testing.T) {
	strategy := &stubStrategy{
		allowAllowed:   true,
		allowRemaining: 2,
	}
	l, err := NewLimiter(strategy, 5, 10, time.Second, time.Minute)
	if err != nil {
		t.Fatalf("NewLimiter error: %v", err)
	}

	allowed, remaining, _, err := l.Check("ip:1.2.3.4", 5)
	if err != nil {
		t.Fatalf("Check error: %v", err)
	}
	if !allowed {
		t.Fatalf("expected allowed")
	}
	if remaining != 2 {
		t.Fatalf("unexpected remaining: %d", remaining)
	}
	if strategy.blockCalls != 0 {
		t.Fatalf("expected no block calls")
	}
	if strategy.lastKey != "ip:1.2.3.4" {
		t.Fatalf("unexpected key: %s", strategy.lastKey)
	}
	if strategy.lastMaxReq != 5 {
		t.Fatalf("unexpected maxReq: %d", strategy.lastMaxReq)
	}
}

func TestLimiterCheck_BlocksWhenNotAllowed(t *testing.T) {
	strategy := &stubStrategy{
		allowAllowed: false,
	}
	l, err := NewLimiter(strategy, 5, 10, time.Second, time.Minute)
	if err != nil {
		t.Fatalf("NewLimiter error: %v", err)
	}

	allowed, _, _, err := l.Check("ip:1.2.3.4", 5)
	if err != nil {
		t.Fatalf("Check error: %v", err)
	}
	if allowed {
		t.Fatalf("expected not allowed")
	}
	if strategy.blockCalls != 1 {
		t.Fatalf("expected block call")
	}
}

func TestLimiterCheck_PropagatesAllowError(t *testing.T) {
	strategy := &stubStrategy{
		allowErr: errors.New("boom"),
	}
	l, err := NewLimiter(strategy, 5, 10, time.Second, time.Minute)
	if err != nil {
		t.Fatalf("NewLimiter error: %v", err)
	}

	if _, _, _, err := l.Check("ip:1.2.3.4", 5); err == nil {
		t.Fatalf("expected error")
	}
	if strategy.blockCalls != 0 {
		t.Fatalf("expected no block calls on error")
	}
}
