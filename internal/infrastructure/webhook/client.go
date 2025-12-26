package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"incident-system/internal/config"
	"incident-system/internal/domain/models"
	"incident-system/pkg/logger"
)

type WebhookClient struct {
    client      *http.Client
    url         string
    maxRetries  int
    retryDelay  time.Duration
    logger      *logger.Logger
}

func NewWebhookClient(cfg *config.Config, logger *logger.Logger) *WebhookClient {
    return &WebhookClient{
        client: &http.Client{
            Timeout: cfg.WebhookTimeout,
        },
        url:        cfg.WebhookURL,
        maxRetries: cfg.WebhookMaxRetries,
        retryDelay: cfg.WebhookRetryDelay,
        logger:     logger,
    }
}

func (w *WebhookClient) Send(payload models.WebhookPayload) error {
    data, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("failed to marshal payload: %w", err)
    }
    
    var lastErr error
    for attempt := 1; attempt <= w.maxRetries; attempt++ {
        w.logger.Info("Sending webhook attempt %d/%d", attempt, w.maxRetries) // Изменено с Infof на Info
        
        req, err := http.NewRequest("POST", w.url, bytes.NewBuffer(data))
        if err != nil {
            lastErr = err
            continue
        }
        
        req.Header.Set("Content-Type", "application/json")
        
        resp, err := w.client.Do(req)
        if err != nil {
            lastErr = err
            w.logger.Error("Webhook attempt %d failed: %v", attempt, err) // Изменено с Errorf на Error
            time.Sleep(w.retryDelay * time.Duration(attempt))
            continue
        }
        
        defer resp.Body.Close()
        
        if resp.StatusCode >= 200 && resp.StatusCode < 300 {
            w.logger.Info("Webhook sent successfully") // Изменено с Infof на Info
            return nil
        }
        
        lastErr = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
        w.logger.Error("Webhook attempt %d failed with status: %d", attempt, resp.StatusCode) // Изменено с Errorf на Error
        time.Sleep(w.retryDelay * time.Duration(attempt))
    }
    
    return fmt.Errorf("failed to send webhook after %d attempts: %w", w.maxRetries, lastErr)
}

func (w *WebhookClient) StartWorker(ctx context.Context, dequeueFunc func(context.Context) (*models.WebhookPayload, error)) {
    go func() {
        for {
            select {
            case <-ctx.Done():
                w.logger.Info("Webhook worker stopped")
                return
            default:
                payload, err := dequeueFunc(ctx)
                if err != nil {
                    w.logger.Error("Failed to dequeue webhook: %v", err) // Изменено с Errorf на Error
                    time.Sleep(time.Second)
                    continue
                }
                
                if payload != nil {
                    if err := w.Send(*payload); err != nil {
                        w.logger.Error("Failed to send webhook: %v", err) // Изменено с Errorf на Error
                    }
                }
            }
        }
    }()
}