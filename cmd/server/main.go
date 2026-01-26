package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/sergioc0sta/limit-barrier/configs"
	"github.com/sergioc0sta/limit-barrier/internal/limiter"
	"github.com/sergioc0sta/limit-barrier/internal/middleware"
	"github.com/sergioc0sta/limit-barrier/internal/storage"
)

func main() {
	v, err := configs.NewConfig()
	if err != nil {
		panic(err)
	}

	store, err := storage.NewStore(v)
	if err != nil {
		panic(err)
	}
	defer store.Close()

	rateLimitDur, err := time.ParseDuration(v.RateLimitDur)
	if err != nil {
		panic(err)
	}
	blockTime := time.Duration(v.BlockTime) * time.Second

	lim, err := limiter.NewLimiter(store, v.IPMaxReq, v.TokenMaxReq, rateLimitDur, blockTime)
	if err != nil {
		panic(err)
	}

	middlewareRateLimiter := middleware.RateLimiter(lim, v.TokenLimits)

	if err := store.Ping(context.Background()); err != nil {
		panic("Storage is not working: " + err.Error())
	} else {
		println("Storage is working.......")
	}

	println("Config loaded:", v.RedisAddr)

	r := gin.Default()
	r.Use(middlewareRateLimiter)
	r.GET("/ping", func(c *gin.Context) {
		ip := c.GetString("client_ip")
		log.Printf("Client IP: %s made request!", ip)
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8080")
}
