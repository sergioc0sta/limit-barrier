package storage

import (
	"fmt"
	"strings"

	"github.com/sergioc0sta/limit-barrier/configs"
	"github.com/sergioc0sta/limit-barrier/internal/storage/redis"
)

func NewStore(cfg *configs.Config) (Store, error) {
	switch strings.ToLower(cfg.StorageDriver) {
	case "redis":
		return redis.NewStore(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB), nil
	case "mysql":
		return nil, fmt.Errorf("storage driver %q not implemented yet", cfg.StorageDriver)
	default:
		return nil, fmt.Errorf("unknown storage driver %q", cfg.StorageDriver)
	}
}
