package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/sergioc0sta/limit-barrier/internal/utils"
)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {

		clientIP, apiKey := utils.GetIpAndToken(c)

		if clientIP == "192.168.0.101" {
			log.Printf("===  Acesso bloqueado para o IP: %s  ===", clientIP)
			c.String(http.StatusTooManyRequests, "You have reached the maximum number of requests or actions allowed within a certain time frame")
			c.Abort()
			return

		}
		log.Printf("[%s] %s %s from IP: %s e API_KEY: %s",
			time.Now().Format("2006-01-02 15:04:05"),
			c.Request.Method,
			c.FullPath(),
			clientIP,
			apiKey)

		c.Next()
	}
}
