package utils

import (
	"github.com/gin-gonic/gin"
)

func GetIpAndToken(c *gin.Context) (string, string) {
	clientIP := c.ClientIP()
	c.Set("client_ip", clientIP)
	apiKey := c.GetHeader("API_KEY")
	c.Set("api_key", apiKey)
	return clientIP, apiKey
}
