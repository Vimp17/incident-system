package services

import (
	"context"
	"fmt"
	"math"
	"time"

	"incident-system/internal/domain/models"
	"incident-system/internal/domain/repositories"
)

type IncidentService struct {
    incidentRepo repositories.IncidentRepository
    cacheRepo    repositories.CacheRepository
    queueRepo    repositories.QueueRepository
}

func NewIncidentService(
    incidentRepo repositories.IncidentRepository,
    cacheRepo repositories.CacheRepository,
    queueRepo repositories.QueueRepository,
) *IncidentService {
    return &IncidentService{
        incidentRepo: incidentRepo,
        cacheRepo:    cacheRepo,
        queueRepo:    queueRepo,
    }
}

func (s *IncidentService) CreateIncident(ctx context.Context, req models.CreateIncidentRequest) (*models.Incident, error) {
    incident := &models.Incident{
        UserID:      req.UserID,
        Latitude:    req.Latitude,
        Longitude:   req.Longitude,
        Title:       req.Title,
        Description: req.Description,
        Severity:    req.Severity,
        Radius:      req.Radius,
        Active:      true,
    }
    
    if err := s.incidentRepo.Create(ctx, incident); err != nil {
        return nil, fmt.Errorf("failed to create incident: %w", err)
    }
    
    // Инвалидируем кеш активных инцидентов
    _ = s.cacheRepo.InvalidateActiveIncidents(ctx)
    
    return incident, nil
}

func (s *IncidentService) GetIncident(ctx context.Context, id int64) (*models.Incident, error) {
    return s.incidentRepo.FindByID(ctx, id)
}

func (s *IncidentService) ListIncidents(ctx context.Context, limit, offset int, activeOnly bool) ([]*models.Incident, int, error) {
    incidents, err := s.incidentRepo.FindAll(ctx, limit, offset, activeOnly)
    if err != nil {
        return nil, 0, err
    }
    
    total, err := s.incidentRepo.CountAll(ctx, activeOnly)
    if err != nil {
        return nil, 0, err
    }
    
    return incidents, total, nil
}

func (s *IncidentService) UpdateIncident(ctx context.Context, id int64, req models.UpdateIncidentRequest) (*models.Incident, error) {
    incident, err := s.incidentRepo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to find incident: %w", err)
    }
    
    if incident == nil {
        return nil, fmt.Errorf("incident not found")
    }
    
    // Обновляем только переданные поля
    if req.Title != nil {
        incident.Title = *req.Title
    }
    if req.Description != nil {
        incident.Description = *req.Description
    }
    if req.Severity != nil {
        incident.Severity = *req.Severity
    }
    if req.Radius != nil {
        incident.Radius = *req.Radius
    }
    if req.Active != nil {
        incident.Active = *req.Active
    }
    
    if err := s.incidentRepo.Update(ctx, incident); err != nil {
        return nil, fmt.Errorf("failed to update incident: %w", err)
    }
    
    // Инвалидируем кеш активных инцидентов
    _ = s.cacheRepo.InvalidateActiveIncidents(ctx)
    
    return incident, nil
}

func (s *IncidentService) DeleteIncident(ctx context.Context, id int64) error {
    if err := s.incidentRepo.Delete(ctx, id); err != nil {
        return fmt.Errorf("failed to delete incident: %w", err)
    }
    
    // Инвалидируем кеш активных инцидентов
    _ = s.cacheRepo.InvalidateActiveIncidents(ctx)
    
    return nil
}

func (s *IncidentService) CheckLocation(ctx context.Context, req models.LocationCheckRequest) (*models.LocationCheckResponse, error) {
    // Сначала пытаемся получить активные инциденты из кеша
    incidents, err := s.cacheRepo.GetActiveIncidents(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get cached incidents: %w", err)
    }
    
    // Если в кеше нет, получаем из базы и кешируем
    if incidents == nil {
        incidents, err = s.incidentRepo.GetActiveIncidents(ctx)
        if err != nil {
            return nil, fmt.Errorf("failed to get active incidents: %w", err)
        }
        
        if err := s.cacheRepo.SetActiveIncidents(ctx, incidents); err != nil {
            // Логируем ошибку, но продолжаем работу
            fmt.Printf("Failed to cache incidents: %v\n", err)
        }
    }
    
    // Фильтруем инциденты по расстоянию
    var nearbyIncidents []models.Incident
    for _, incident := range incidents {
        distance := calculateDistance(
            req.Latitude, req.Longitude,
            incident.Latitude, incident.Longitude,
        )
        
        // Переводим радиус из метров в километры
        if distance <= incident.Radius/1000.0 {
            nearbyIncidents = append(nearbyIncidents, *incident)
        }
    }
    
    hasAlert := len(nearbyIncidents) > 0
    
    // Сохраняем факт проверки
    check := &models.LocationCheck{
        UserID:    req.UserID,
        Latitude:  req.Latitude,
        Longitude: req.Longitude,
        Timestamp: time.Now(),
        HasAlert:  hasAlert,
    }
    
    if err := s.incidentRepo.SaveLocationCheck(ctx, check); err != nil {
        // Логируем ошибку, но не прерываем выполнение
        fmt.Printf("Failed to save location check: %v\n", err)
    }
    
    // Если есть опасные зоны, ставим задачу на отправку вебхука
    if hasAlert {
        go s.enqueueWebhook(ctx, req, nearbyIncidents)
    }
    
    return &models.LocationCheckResponse{
        Incidents: nearbyIncidents,
        HasAlert:  hasAlert,
    }, nil
}

func (s *IncidentService) GetStats(ctx context.Context, minutes int) ([]models.IncidentStats, error) {
    stats, err := s.incidentRepo.GetStats(ctx, minutes)
    if err != nil {
        return nil, fmt.Errorf("failed to get stats: %w", err)
    }
    
    result := make([]models.IncidentStats, len(stats))
    for i, stat := range stats {
        result[i] = *stat
    }
    
    return result, nil
}

func (s *IncidentService) enqueueWebhook(ctx context.Context, req models.LocationCheckRequest, incidents []models.Incident) {
    // Создаем укороченную версию инцидентов для вебхука
    var shortIncidents []models.IncidentShort
    for _, incident := range incidents {
        distance := calculateDistance(
            req.Latitude, req.Longitude,
            incident.Latitude, incident.Longitude,
        )
        
        shortIncidents = append(shortIncidents, models.IncidentShort{
            ID:       incident.ID,
            Title:    incident.Title,
            Severity: incident.Severity,
            Distance: distance * 1000, // переводим в метры
        })
    }
    
    payload := models.WebhookPayload{
        EventType: "location_alert",
        UserID:    req.UserID,
        Latitude:  req.Latitude,
        Longitude: req.Longitude,
        Incidents: shortIncidents,
        Timestamp: time.Now(),
    }
    
    if err := s.queueRepo.EnqueueWebhook(ctx, payload); err != nil {
        fmt.Printf("Failed to enqueue webhook: %v\n", err)
    }
}

// calculateDistance вычисляет расстояние между двумя точками по формуле гаверсинусов
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
    const R = 6371 // радиус Земли в километрах
    
    dLat := (lat2 - lat1) * math.Pi / 180
    dLon := (lon2 - lon1) * math.Pi / 180
    
    a := math.Sin(dLat/2)*math.Sin(dLat/2) +
        math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
            math.Sin(dLon/2)*math.Sin(dLon/2)
    
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    
    return R * c
}