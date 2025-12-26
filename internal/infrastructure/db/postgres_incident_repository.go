package db

import (
	"context"
	"database/sql"
	"time"

	"incident-system/internal/domain/models"
	"incident-system/internal/domain/repositories"
)

type postgresIncidentRepository struct {
    db *sql.DB
}

func NewPostgresIncidentRepository(db *sql.DB) repositories.IncidentRepository {
    return &postgresIncidentRepository{db: db}
}

func (r *postgresIncidentRepository) Create(ctx context.Context, incident *models.Incident) error {
    query := `
        INSERT INTO incidents (
            user_id, latitude, longitude, title, description, 
            severity, radius, active, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id
    `
    
    now := time.Now()
    incident.CreatedAt = now
    incident.UpdatedAt = now
    
    err := r.db.QueryRowContext(ctx, query,
        incident.UserID,
        incident.Latitude,
        incident.Longitude,
        incident.Title,
        incident.Description,
        incident.Severity,
        incident.Radius,
        incident.Active,
        incident.CreatedAt,
        incident.UpdatedAt,
    ).Scan(&incident.ID)
    
    return err
}

func (r *postgresIncidentRepository) FindByID(ctx context.Context, id int64) (*models.Incident, error) {
    query := `
        SELECT id, user_id, latitude, longitude, title, description,
               severity, radius, active, created_at, updated_at
        FROM incidents
        WHERE id = $1
    `
    
    var incident models.Incident
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &incident.ID,
        &incident.UserID,
        &incident.Latitude,
        &incident.Longitude,
        &incident.Title,
        &incident.Description,
        &incident.Severity,
        &incident.Radius,
        &incident.Active,
        &incident.CreatedAt,
        &incident.UpdatedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    
    return &incident, err
}

func (r *postgresIncidentRepository) FindAll(ctx context.Context, limit, offset int, activeOnly bool) ([]*models.Incident, error) {
    query := `
        SELECT id, user_id, latitude, longitude, title, description,
               severity, radius, active, created_at, updated_at
        FROM incidents
        WHERE ($1 = false OR active = $1)
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `
    
    rows, err := r.db.QueryContext(ctx, query, activeOnly, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var incidents []*models.Incident
    for rows.Next() {
        var incident models.Incident
        if err := rows.Scan(
            &incident.ID,
            &incident.UserID,
            &incident.Latitude,
            &incident.Longitude,
            &incident.Title,
            &incident.Description,
            &incident.Severity,
            &incident.Radius,
            &incident.Active,
            &incident.CreatedAt,
            &incident.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        incidents = append(incidents, &incident)
    }
    
    return incidents, nil
}

func (r *postgresIncidentRepository) Update(ctx context.Context, incident *models.Incident) error {
    query := `
        UPDATE incidents 
        SET title = $1, description = $2, severity = $3, 
            radius = $4, active = $5, updated_at = $6
        WHERE id = $7
    `
    
    incident.UpdatedAt = time.Now()
    _, err := r.db.ExecContext(ctx, query,
        incident.Title,
        incident.Description,
        incident.Severity,
        incident.Radius,
        incident.Active,
        incident.UpdatedAt,
        incident.ID,
    )
    
    return err
}

func (r *postgresIncidentRepository) Delete(ctx context.Context, id int64) error {
    // Деактивация вместо удаления
    query := `UPDATE incidents SET active = false, updated_at = $1 WHERE id = $2`
    _, err := r.db.ExecContext(ctx, query, time.Now(), id)
    return err
}

func (r *postgresIncidentRepository) FindNearLocation(ctx context.Context, lat, lng float64, radiusKm float64) ([]*models.Incident, error) {
    // Используем формулу гаверсинусов для расчета расстояния
    query := `
        SELECT id, user_id, latitude, longitude, title, description,
               severity, radius, active, created_at, updated_at,
               (6371 * acos(
                   cos(radians($1)) * cos(radians(latitude)) * 
                   cos(radians(longitude) - radians($2)) + 
                   sin(radians($1)) * sin(radians(latitude))
               )) as distance
        FROM incidents
        WHERE active = true
        HAVING distance <= $3
        ORDER BY distance
    `
    
    rows, err := r.db.QueryContext(ctx, query, lat, lng, radiusKm)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var incidents []*models.Incident
    for rows.Next() {
        var incident models.Incident
        var distance float64
        
        if err := rows.Scan(
            &incident.ID,
            &incident.UserID,
            &incident.Latitude,
            &incident.Longitude,
            &incident.Title,
            &incident.Description,
            &incident.Severity,
            &incident.Radius,
            &incident.Active,
            &incident.CreatedAt,
            &incident.UpdatedAt,
            &distance,
        ); err != nil {
            return nil, err
        }
        incidents = append(incidents, &incident)
    }
    
    return incidents, nil
}

func (r *postgresIncidentRepository) SaveLocationCheck(ctx context.Context, check *models.LocationCheck) error {
    query := `
        INSERT INTO location_checks (user_id, latitude, longitude, timestamp, has_alert, incident_id)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `
    
    return r.db.QueryRowContext(ctx, query,
        check.UserID,
        check.Latitude,
        check.Longitude,
        check.Timestamp,
        check.HasAlert,
        check.IncidentID, // Может быть NULL
    ).Scan(&check.ID)
}

func (r *postgresIncidentRepository) GetStats(ctx context.Context, minutes int) ([]*models.IncidentStats, error) {
    // Используем COALESCE для обработки NULL значений
    query := `
        SELECT 
            COALESCE(incident_id, 0) as zone_id,
            COUNT(DISTINCT user_id) as user_count
        FROM location_checks
        WHERE timestamp >= NOW() - ($1 || ' minutes')::INTERVAL
        GROUP BY COALESCE(incident_id, 0)
    `
    
    rows, err := r.db.QueryContext(ctx, query, minutes)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var stats []*models.IncidentStats
    for rows.Next() {
        var stat models.IncidentStats
        if err := rows.Scan(&stat.ZoneID, &stat.UserCount); err != nil {
            return nil, err
        }
        // Для zone_id = 0 (NULL значения) устанавливаем nil
        if stat.ZoneID != nil && *stat.ZoneID == 0 {
            stat.ZoneID = nil
        }
        stats = append(stats, &stat)
    }
    
    return stats, nil
}

func (r *postgresIncidentRepository) GetActiveIncidents(ctx context.Context) ([]*models.Incident, error) {
    query := `
        SELECT id, user_id, latitude, longitude, title, description,
               severity, radius, active, created_at, updated_at
        FROM incidents
        WHERE active = true
        ORDER BY id
    `
    
    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var incidents []*models.Incident
    for rows.Next() {
        var incident models.Incident
        if err := rows.Scan(
            &incident.ID,
            &incident.UserID,
            &incident.Latitude,
            &incident.Longitude,
            &incident.Title,
            &incident.Description,
            &incident.Severity,
            &incident.Radius,
            &incident.Active,
            &incident.CreatedAt,
            &incident.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        incidents = append(incidents, &incident)
    }
    
    return incidents, nil
}

func (r *postgresIncidentRepository) CountAll(ctx context.Context, activeOnly bool) (int, error) {
    query := `SELECT COUNT(*) FROM incidents WHERE ($1 = false OR active = $1)`
    
    var count int
    err := r.db.QueryRowContext(ctx, query, activeOnly).Scan(&count)
    return count, err
}