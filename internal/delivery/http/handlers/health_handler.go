package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
    db    *sql.DB
    redis *redis.Client
}

func NewHealthHandler(db *sql.DB, redis *redis.Client) *HealthHandler {
    return &HealthHandler{
        db:    db,
        redis: redis,
    }
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
    // Проверка базы данных
    if err := h.db.Ping(); err != nil {
        c.JSON(http.StatusServiceUnavailable, gin.H{
            "status": "unhealthy",
            "error":  "database connection failed",
        })
        return
    }
    
    // Проверка Redis
    if err := h.redis.Ping(c.Request.Context()).Err(); err != nil {
        c.JSON(http.StatusServiceUnavailable, gin.H{
            "status": "unhealthy",
            "error":  "redis connection failed",
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "status": "healthy",
        "services": gin.H{
            "database": "connected",
            "redis":    "connected",
        },
    })
}