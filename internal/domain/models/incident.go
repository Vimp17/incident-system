package models

import (
	"time"
)

type Incident struct {
    ID          int64     `json:"id" db:"id"`
    UserID      string    `json:"user_id" db:"user_id"`
    Latitude    float64   `json:"latitude" db:"latitude"`
    Longitude   float64   `json:"longitude" db:"longitude"`
    Title       string    `json:"title" db:"title"`
    Description string    `json:"description" db:"description"`
    Severity    string    `json:"severity" db:"severity"` // low, medium, high
    Radius      float64   `json:"radius" db:"radius"` // в метрах
    Active      bool      `json:"active" db:"active"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type IncidentStats struct {
    ZoneID    *int64 `json:"zone_id" db:"zone_id"` 
    UserCount int64  `json:"user_count" db:"user_count"`
}

type CreateIncidentRequest struct {
    UserID      string  `json:"user_id" validate:"required"`
    Latitude    float64 `json:"latitude" validate:"required,latitude"`
    Longitude   float64 `json:"longitude" validate:"required,longitude"`
    Title       string  `json:"title" validate:"required,min=3,max=255"`
    Description string  `json:"description" validate:"max=1000"`
    Severity    string  `json:"severity" validate:"required,oneof=low medium high"`
    Radius      float64 `json:"radius" validate:"required,min=10,max=5000"`
}

type UpdateIncidentRequest struct {
    Title       *string  `json:"title" validate:"omitempty,min=3,max=255"`
    Description *string  `json:"description" validate:"omitempty,max=1000"`
    Severity    *string  `json:"severity" validate:"omitempty,oneof=low medium high"`
    Radius      *float64 `json:"radius" validate:"omitempty,min=10,max=5000"`
    Active      *bool    `json:"active"`
}