package services

import (
	"context"
	"time"

	"incident-system/internal/domain/models"
	"incident-system/internal/domain/repositories"
	"incident-system/pkg/logger"
)

type WebhookService struct {
    queueRepo repositories.QueueRepository
    logger    *logger.Logger
}

func NewWebhookService(queueRepo repositories.QueueRepository, logger *logger.Logger) *WebhookService {
    return &WebhookService{
        queueRepo: queueRepo,
        logger:    logger,
    }
}

func (s *WebhookService) EnqueueWebhook(ctx context.Context, payload models.WebhookPayload) error {
    return s.queueRepo.EnqueueWebhook(ctx, payload)
}

func (s *WebhookService) StartWorker(ctx context.Context, processFunc func(models.WebhookPayload) error) {
    go func() {
        for {
            select {
            case <-ctx.Done():
                s.logger.Info("Webhook worker stopped")
                return
            default:
                payload, err := s.queueRepo.DequeueWebhook(ctx)
                if err != nil {
                    s.logger.Error("Failed to dequeue webhook: %v", err) // Изменено с Errorf на Error
                    time.Sleep(time.Second)
                    continue
                }
                
                if payload != nil {
                    if err := processFunc(*payload); err != nil {
                        s.logger.Error("Failed to process webhook: %v", err) // Изменено с Errorf на Error
                    }
                }
            }
        }
    }()
}