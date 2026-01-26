package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	blockPrefix = "rl:block:"
	countPrefix = "rl:count:"
)

func (s *Store) Allow(key string, maxReq int, duration time.Duration) (bool, int, time.Time, error) {
	if key == "" {
		return false, 0, time.Time{}, errors.New("key is required")
	}
	if maxReq <= 0 {
		return false, 0, time.Time{}, errors.New("maxReq must be > 0")
	}
	if duration <= 0 {
		return false, 0, time.Time{}, errors.New("duration must be > 0")
	}

	blockKey := blockPrefix + key
	countKey := countPrefix + key

	ctx := context.Background()
	blocked, err := s.client.Exists(ctx, blockKey).Result()
	if err != nil {
		return false, 0, time.Time{}, fmt.Errorf("check block: %w", err)
	}
	if blocked == 1 {
		ttl, err := s.client.PTTL(ctx, blockKey).Result()
		if err != nil && err != redis.Nil {
			return false, 0, time.Time{}, fmt.Errorf("block ttl: %w", err)
		}
		var reset time.Time
		if ttl > 0 {
			reset = time.Now().Add(ttl)
		}
		return false, 0, reset, nil
	}

	pipe := s.client.TxPipeline()
	incrCmd := pipe.Incr(ctx, countKey)
	ttlCmd := pipe.PTTL(ctx, countKey)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return false, 0, time.Time{}, fmt.Errorf("incr/ttl: %w", err)
	}

	count, err := incrCmd.Result()
	if err != nil {
		return false, 0, time.Time{}, fmt.Errorf("count: %w", err)
	}
	ttl, err := ttlCmd.Result()
	if err != nil && err != redis.Nil {
		return false, 0, time.Time{}, fmt.Errorf("ttl: %w", err)
	}

	if ttl <= 0 {
		if err := s.client.PExpire(ctx, countKey, duration).Err(); err != nil {
			return false, 0, time.Time{}, fmt.Errorf("set ttl: %w", err)
		}
		ttl = duration
	}

	var reset time.Time
	if ttl > 0 {
		reset = time.Now().Add(ttl)
	}

	if int(count) > maxReq {
		return false, 0, reset, nil
	}

	remaining := maxReq - int(count)
	if remaining < 0 {
		remaining = 0
	}

	return true, remaining, reset, nil
}

func (s *Store) Block(key string, duration time.Duration) error {
	if key == "" {
		return errors.New("key is required")
	}
	if duration <= 0 {
		return errors.New("duration must be > 0")
	}

	ctx := context.Background()
	blockKey := blockPrefix + key
	return s.client.Set(ctx, blockKey, "1", duration).Err()
}

func (s *Store) Unblock(key string) error {
	if key == "" {
		return errors.New("key is required")
	}

	ctx := context.Background()
	blockKey := blockPrefix + key
	return s.client.Del(ctx, blockKey).Err()
}
