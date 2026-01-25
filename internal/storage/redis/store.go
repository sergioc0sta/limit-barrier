package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
}

func NewStore(addr, password string, db int) *Store {
	return &Store{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}
}

func (s *Store) Ping(ctx context.Context) error {
	return s.client.Ping(ctx).Err()
}

func (s *Store) Set(ctx context.Context, key, value string) error {
	return s.client.Set(ctx, key, value, 0).Err()
}

func (s *Store) Get(ctx context.Context, key string) (string, error) {
	return s.client.Get(ctx, key).Result()
}

func (s *Store) Close() error {
	return s.client.Close()
}
