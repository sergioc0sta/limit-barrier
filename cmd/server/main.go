package main

import (
	"github.com/gin-gonic/gin"

	"github.com/sergioc0sta/limit-barrier/configs"
)

func main() {
	v, err := configs.NewConfig()
	if err != nil {
		panic(err)
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
