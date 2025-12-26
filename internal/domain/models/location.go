package models

import (
	"time"
)

type LocationCheck struct {
    ID         int64     `json:"id" db:"id"`
    UserID     string    `json:"user_id" db:"user_id"`
    Latitude   float64   `json:"latitude" db:"latitude"`
    Longitude  float64   `json:"longitude" db:"longitude"`
    Timestamp  time.Time `json:"timestamp" db:"timestamp"`
    HasAlert   bool      `json:"has_alert" db:"has_alert"`
    IncidentID *int64    `json:"incident_id,omitempty" db:"incident_id"` // Изменено на указатель
}

type LocationCheckRequest struct {
    UserID    string  `json:"user_id" validate:"required"`
    Latitude  float64 `json:"latitude" validate:"required,latitude"`
    Longitude float64 `json:"longitude" validate:"required,longitude"`
}

type LocationCheckResponse struct {
    Incidents []Incident `json:"incidents"`
    HasAlert  bool       `json:"has_alert"`
}

type WebhookPayload struct {
    EventType string          `json:"event_type"` // "location_alert"
    UserID    string          `json:"user_id"`
    Latitude  float64         `json:"latitude"`
    Longitude float64         `json:"longitude"`
    Incidents []IncidentShort `json:"incidents"`
    Timestamp time.Time       `json:"timestamp"`
}

type IncidentShort struct {
    ID       int64   `json:"id"`
    Title    string  `json:"title"`
    Severity string  `json:"severity"`
    Distance float64 `json:"distance"` // расстояние в метрах
}