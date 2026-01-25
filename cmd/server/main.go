package main

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/sergioc0sta/limit-barrier/configs"
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

	if err := store.Ping(context.Background()); err != nil {
		panic("Storage is not working: " + err.Error())
	} else {
		println("Storage is working.......")
	}

	println("Config loaded:", v.RedisAddr)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8080")
}
