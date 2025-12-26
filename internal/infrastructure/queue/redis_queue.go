package queue

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

type redisQueueRepository struct {
    client *redis.Client
    queue  string
}

func NewRedisQueueRepository(cfg *config.Config) (repositories.QueueRepository, error) {
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
    
    return &redisQueueRepository{
        client: client,
        queue:  "webhook_queue",
    }, nil
}

func (r *redisQueueRepository) EnqueueWebhook(ctx context.Context, payload models.WebhookPayload) error {
    data, err := json.Marshal(payload)
    if err != nil {
        return err
    }
    
    return r.client.LPush(ctx, r.queue, data).Err()
}

func (r *redisQueueRepository) DequeueWebhook(ctx context.Context) (*models.WebhookPayload, error) {
    result, err := r.client.BRPop(ctx, 0, r.queue).Result()
    if err != nil {
        return nil, err
    }
    
    if len(result) < 2 {
        return nil, fmt.Errorf("invalid queue result")
    }
    
    var payload models.WebhookPayload
    if err := json.Unmarshal([]byte(result[1]), &payload); err != nil {
        return nil, err
    }
    
    return &payload, nil
}