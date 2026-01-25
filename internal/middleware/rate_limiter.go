package middleware

import (
	"log"
	"time"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {

		clientIP := c.ClientIP()
		c.Set("client_ip", clientIP)
		apiKey := c.GetHeader("API_KEY")
		c.Set("api_key", apiKey)

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
