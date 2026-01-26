package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/sergioc0sta/limit-barrier/internal/limiter"
)

type recordStrategy struct {
	allowed    bool
	allowErr   error
	lastKey    string
	lastMaxReq int
	lastWindow time.Duration
	blockCalls int
	blockErr   error
}

func (r *recordStrategy) Allow(key string, maxReq int, duration time.Duration) (bool, int, time.Time, error) {
	r.lastKey = key
	r.lastMaxReq = maxReq
	r.lastWindow = duration
	return r.allowed, 0, time.Time{}, r.allowErr
}

func (r *recordStrategy) Block(key string, duration time.Duration) error {
	r.blockCalls++
	return r.blockErr
}

func (r *recordStrategy) Unblock(key string) error {
	return nil
}

func newRouter(l *limiter.Limiter, tokenLimits map[string]int) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RateLimiter(l, tokenLimits))
	r.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	return r
}

func TestMiddleware_UsesIPWhenNoToken(t *testing.T) {
	strategy := &recordStrategy{allowed: true}
	l, err := limiter.NewLimiter(strategy, 5, 10, time.Second, time.Minute)
	if err != nil {
		t.Fatalf("NewLimiter error: %v", err)
	}

	r := newRouter(l, map[string]int{"TOKEN_BASIC": 100})
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.RemoteAddr = "203.0.113.5:1234"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", w.Code)
	}
	if strategy.lastKey != "ip:203.0.113.5" {
		t.Fatalf("unexpected key: %s", strategy.lastKey)
	}
	if strategy.lastMaxReq != 5 {
		t.Fatalf("unexpected maxReq: %d", strategy.lastMaxReq)
	}
}

func TestMiddleware_TokenFallbacksToIPWhenNotInJSON(t *testing.T) {
	strategy := &recordStrategy{allowed: true}
	l, err := limiter.NewLimiter(strategy, 5, 10, time.Second, time.Minute)
	if err != nil {
		t.Fatalf("NewLimiter error: %v", err)
	}

	r := newRouter(l, map[string]int{"TOKEN_BASIC": 100})
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("API_KEY", "TOKEN_UNKNOWN")
	req.RemoteAddr = "203.0.113.6:1234"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", w.Code)
	}
	if strategy.lastKey != "ip:203.0.113.6" {
		t.Fatalf("unexpected key: %s", strategy.lastKey)
	}
	if strategy.lastMaxReq != 5 {
		t.Fatalf("unexpected maxReq: %d", strategy.lastMaxReq)
	}
}

func TestMiddleware_UsesTokenWhenInJSON(t *testing.T) {
	strategy := &recordStrategy{allowed: true}
	l, err := limiter.NewLimiter(strategy, 5, 10, time.Second, time.Minute)
	if err != nil {
		t.Fatalf("NewLimiter error: %v", err)
	}

	r := newRouter(l, map[string]int{"TOKEN_BASIC": 100})
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("API_KEY", "TOKEN_BASIC")
	req.RemoteAddr = "203.0.113.7:1234"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", w.Code)
	}
	if strategy.lastKey != "token:TOKEN_BASIC" {
		t.Fatalf("unexpected key: %s", strategy.lastKey)
	}
	if strategy.lastMaxReq != 100 {
		t.Fatalf("unexpected maxReq: %d", strategy.lastMaxReq)
	}
}

func TestMiddleware_Returns429WhenLimited(t *testing.T) {
	strategy := &recordStrategy{allowed: false}
	l, err := limiter.NewLimiter(strategy, 5, 10, time.Second, time.Minute)
	if err != nil {
		t.Fatalf("NewLimiter error: %v", err)
	}

	r := newRouter(l, map[string]int{})
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.RemoteAddr = "203.0.113.8:1234"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("unexpected status: %d", w.Code)
	}
	if w.Body.String() != limitExceededMessage {
		t.Fatalf("unexpected body: %q", w.Body.String())
	}
}
