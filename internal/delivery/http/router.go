package http

import (
	"context"
	"database/sql"

	"incident-system/internal/config"
	"incident-system/internal/delivery/http/handlers"
	"incident-system/internal/delivery/http/middleware"
	"incident-system/internal/infrastructure/cache"
	"incident-system/internal/infrastructure/db"
	"incident-system/internal/infrastructure/queue"
	"incident-system/internal/infrastructure/webhook"
	"incident-system/internal/usecase/services"
	"incident-system/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupRouter(
    ctx context.Context, // Добавили контекст как первый параметр
    cfg *config.Config,
    postgresDB *sql.DB,
    redisClient *redis.Client,
    logger *logger.Logger,
) *gin.Engine {
    if cfg.Environment == "production" {
        gin.SetMode(gin.ReleaseMode)
    }
    
    router := gin.Default()
    
    // Инициализация репозиториев
    incidentRepo := db.NewPostgresIncidentRepository(postgresDB)
    cacheRepo, _ := cache.NewRedisCacheRepository(cfg)
    queueRepo, _ := queue.NewRedisQueueRepository(cfg)
    
    // Инициализация сервисов
    incidentService := services.NewIncidentService(incidentRepo, cacheRepo, queueRepo)
    
    // Инициализация обработчиков
    incidentHandler := handlers.NewIncidentHandler(incidentService)
    locationHandler := handlers.NewLocationHandler(incidentService)
    healthHandler := handlers.NewHealthHandler(postgresDB, redisClient)
    
    // Инициализация вебхук клиента
    webhookClient := webhook.NewWebhookClient(cfg, logger)
    
    // Запуск воркера для отправки вебхуков
    webhookClient.StartWorker(ctx, queueRepo.DequeueWebhook)
    
    // Public routes
    public := router.Group("/api/v1")
    {
        public.POST("/location/check", locationHandler.CheckLocation)
        public.GET("/system/health", healthHandler.HealthCheck)
    }
    
    // Protected routes (требуют API key)
    protected := router.Group("/api/v1")
    protected.Use(middleware.APIKeyAuth(cfg))
    {
        // CRUD для инцидентов
        incidents := protected.Group("/incidents")
        {
            incidents.POST("", incidentHandler.CreateIncident)
            incidents.GET("", incidentHandler.ListIncidents)
            incidents.GET("/:id", incidentHandler.GetIncident)
            incidents.PUT("/:id", incidentHandler.UpdateIncident)
            incidents.DELETE("/:id", incidentHandler.DeleteIncident)
        }
        
        // Статистика
        protected.GET("/incidents/stats", incidentHandler.GetStats)
    }
    
    return router
}