package repositories

import (
	"context"
	"incident-system/internal/domain/models"
)

type IncidentRepository interface {
    // CRUD операции
    Create(ctx context.Context, incident *models.Incident) error
    FindByID(ctx context.Context, id int64) (*models.Incident, error)
    FindAll(ctx context.Context, limit, offset int, activeOnly bool) ([]*models.Incident, error)
    Update(ctx context.Context, incident *models.Incident) error
    Delete(ctx context.Context, id int64) error
    
    // Специфичные операции
    FindNearLocation(ctx context.Context, lat, lng float64, radiusKm float64) ([]*models.Incident, error)
    SaveLocationCheck(ctx context.Context, check *models.LocationCheck) error
    GetStats(ctx context.Context, minutes int) ([]*models.IncidentStats, error)
    GetActiveIncidents(ctx context.Context) ([]*models.Incident, error)
    CountAll(ctx context.Context, activeOnly bool) (int, error)
}

type CacheRepository interface {
    GetActiveIncidents(ctx context.Context) ([]*models.Incident, error)
    SetActiveIncidents(ctx context.Context, incidents []*models.Incident) error
    InvalidateActiveIncidents(ctx context.Context) error
}

type QueueRepository interface {
    EnqueueWebhook(ctx context.Context, payload models.WebhookPayload) error
    DequeueWebhook(ctx context.Context) (*models.WebhookPayload, error)
}