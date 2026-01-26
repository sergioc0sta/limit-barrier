package redis

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
)

func newTestStore(t *testing.T) (*Store, *miniredis.Miniredis) {
	t.Helper()

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("miniredis run: %v", err)
	}
	store := NewStore(mr.Addr(), "", 0)
	return store, mr
}

func TestAllow_CountsAndLimits(t *testing.T) {
	store, mr := newTestStore(t)
	defer mr.Close()

	key := "ip:1.2.3.4"
	window := 2 * time.Second

	allowed, remaining, _, err := store.Allow(key, 2, window)
	if err != nil {
		t.Fatalf("Allow error: %v", err)
	}
	if !allowed || remaining != 1 {
		t.Fatalf("unexpected result: allowed=%v remaining=%d", allowed, remaining)
	}
	if ttl := mr.TTL(countPrefix + key); ttl <= 0 {
		t.Fatalf("expected ttl to be set")
	}

	allowed, remaining, _, err = store.Allow(key, 2, window)
	if err != nil {
		t.Fatalf("Allow error: %v", err)
	}
	if !allowed || remaining != 0 {
		t.Fatalf("unexpected result: allowed=%v remaining=%d", allowed, remaining)
	}

	allowed, _, _, err = store.Allow(key, 2, window)
	if err != nil {
		t.Fatalf("Allow error: %v", err)
	}
	if allowed {
		t.Fatalf("expected not allowed")
	}
}

func TestAllow_WhenBlocked(t *testing.T) {
	store, mr := newTestStore(t)
	defer mr.Close()

	key := "ip:1.2.3.4"
	if err := store.Block(key, time.Second); err != nil {
		t.Fatalf("Block error: %v", err)
	}

	allowed, _, _, err := store.Allow(key, 2, time.Second)
	if err != nil {
		t.Fatalf("Allow error: %v", err)
	}
	if allowed {
		t.Fatalf("expected blocked to deny")
	}
}

func TestBlockUnblock(t *testing.T) {
	store, mr := newTestStore(t)
	defer mr.Close()

	key := "ip:1.2.3.4"
	if err := store.Block(key, time.Second); err != nil {
		t.Fatalf("Block error: %v", err)
	}
	if !mr.Exists(blockPrefix + key) {
		t.Fatalf("expected block key to exist")
	}

	if err := store.Unblock(key); err != nil {
		t.Fatalf("Unblock error: %v", err)
	}
	if mr.Exists(blockPrefix + key) {
		t.Fatalf("expected block key to be removed")
	}
}
