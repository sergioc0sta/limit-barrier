package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	redisClient *redis.Client
	ctx context.Context
}

func NewClient(addr string, password string, db int) *Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		panic("Redis connection failed: " + err.Error())
	}

	return &Client{redisClient: redisClient, ctx: ctx}
}


func (client *Client) Set(key, value string) error {
	return client.redisClient.Set(client.ctx, key, value, 0).Err()
}

func (client *Client) Get(key string) (string, error) {
	return client.redisClient.Get(client.ctx, key).Result()
}

func (client *Client) Close() error {
	return client.redisClient.Close()
}

func (client *Client) IsWorking() error {
	return client.redisClient.Ping(client.ctx).Err()
}
