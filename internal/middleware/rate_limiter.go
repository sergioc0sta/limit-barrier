package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/sergioc0sta/limit-barrier/internal/limiter"
	"github.com/sergioc0sta/limit-barrier/internal/utils"
)

const limitExceededMessage = "you have reached the maximum number of requests or actions allowed within a certain time frame"

func RateLimiter(l *limiter.Limiter, tokenLimits map[string]int) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP, apiKey := utils.GetIpAndToken(c)

		log.Printf("[%s] %s %s from IP: %s e API_KEY: %s",
			time.Now().Format("2006-01-02 15:04:05"),
			c.Request.Method,
			c.FullPath(),
			clientIP,
			apiKey)

		var key string
		var maxReq int
		if apiKey != "" {
			if limit, ok := tokenLimits[apiKey]; ok {
				key = "token:" + apiKey
				maxReq = limit
			} else {
				key = "ip:" + clientIP
				maxReq = l.IPMaxReq()
			}
		} else {
			key = "ip:" + clientIP
			maxReq = l.IPMaxReq()
		}

		allowed, _, _, err := l.Check(key, maxReq)
		if err != nil {
			log.Printf("rate limiter error: %v", err)
			c.String(http.StatusInternalServerError, "internal server error")
			c.Abort()
			return
		}
		if !allowed {
			c.String(http.StatusTooManyRequests, limitExceededMessage)
			c.Abort()
			return
		}

		c.Next()
	}
}
