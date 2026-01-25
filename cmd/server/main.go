package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/sergioc0sta/limit-barrier/configs"
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

	middlewareRateLimiter := middleware.RateLimiter()

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
