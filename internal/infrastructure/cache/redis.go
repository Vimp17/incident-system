package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"incident-system/internal/config"
	"incident-system/internal/domain/models"
	"incident-system/internal/domain/repositories"

	"github.com/redis/go-redis/v9"
)

type redisCacheRepository struct {
    client *redis.Client
    ttl    time.Duration
}

func NewRedisCacheRepository(cfg *config.Config) (repositories.CacheRepository, error) {
    client := redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
        Password: cfg.RedisPassword,
        DB:       cfg.RedisDB,
    })
    
    // Проверка подключения
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := client.Ping(ctx).Err(); err != nil {
        return nil, fmt.Errorf("failed to connect to Redis: %w", err)
    }
    
    return &redisCacheRepository{
        client: client,
        ttl:    time.Duration(cfg.CacheTTLMinutes) * time.Minute,
    }, nil
}

func (r *redisCacheRepository) GetActiveIncidents(ctx context.Context) ([]*models.Incident, error) {
    key := "active_incidents"
    
    data, err := r.client.Get(ctx, key).Result()
    if err == redis.Nil {
        return nil, nil // Ключ не найден - это не ошибка
    }
    if err != nil {
        return nil, err
    }
    
    var incidents []*models.Incident
    if err := json.Unmarshal([]byte(data), &incidents); err != nil {
        return nil, err
    }
    
    return incidents, nil
}

func (r *redisCacheRepository) SetActiveIncidents(ctx context.Context, incidents []*models.Incident) error {
    key := "active_incidents"
    
    data, err := json.Marshal(incidents)
    if err != nil {
        return err
    }
    
    return r.client.Set(ctx, key, data, r.ttl).Err()
}

func (r *redisCacheRepository) InvalidateActiveIncidents(ctx context.Context) error {
    key := "active_incidents"
    return r.client.Del(ctx, key).Err()
}