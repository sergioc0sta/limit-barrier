package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	StorageDriver string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	IPMaxReq      int
	TokenMaxReq   int
	BlockTime     int
	RateLimitDur  string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, err
	}

	ipMaxReq, err := strconv.Atoi(os.Getenv("IP_MAX_REQ"))
	if err != nil {
		return nil, err
	}

	tokenMaxReq, err := strconv.Atoi(os.Getenv("TOKEN_MAX_REQ"))
	if err != nil {
		return nil, err
	}

	blockTime, err := strconv.Atoi(os.Getenv("BLOCK_TIME"))
	if err != nil {
		return nil, err
	}

	return &Config{
		StorageDriver: os.Getenv("STORAGE_DRIVER"),
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       redisDB,
		IPMaxReq:      ipMaxReq,
		TokenMaxReq:   tokenMaxReq,
		BlockTime:     blockTime,
		RateLimitDur:  os.Getenv("RATE_LIMIT_DUR"),
	}, nil
}

