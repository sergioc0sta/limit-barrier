package main

import (
	"github.com/gin-gonic/gin"

	"github.com/sergioc0sta/limit-barrier/configs"
	"github.com/sergioc0sta/limit-barrier/internal/infra/redis"
)

func main() {
	v, err := configs.NewConfig()
	if err != nil {
		panic(err)
	}

	clientRedis := redis.NewClient(v.RedisAddr, v.RedisPassword, v.RedisDB)
	defer clientRedis.Close()

	if err := clientRedis.IsWorking(); err != nil {
		panic("Redis is not working: " + err.Error())
	} else {
		println("Redis is working.......")
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
