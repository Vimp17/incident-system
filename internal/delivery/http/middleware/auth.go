package middleware

import (
	"net/http"

	"incident-system/internal/config"

	"github.com/gin-gonic/gin"
)

func APIKeyAuth(cfg *config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        apiKey := c.GetHeader("X-API-Key")
        
        if apiKey == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
            c.Abort()
            return
        }
        
        if apiKey != cfg.APIKeyOperator {
            c.JSON(http.StatusForbidden, gin.H{"error": "Invalid API key"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}