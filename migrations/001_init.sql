-- Таблица инцидентов
CREATE TABLE incidents (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    severity VARCHAR(50) NOT NULL CHECK (severity IN ('low', 'medium', 'high')),
    radius DECIMAL(10, 2) NOT NULL, -- в метрах
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Таблица проверок локаций
CREATE TABLE location_checks (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    has_alert BOOLEAN NOT NULL,
    incident_id BIGINT REFERENCES incidents(id) ON DELETE SET NULL
);

-- Индексы для быстрого поиска
CREATE INDEX idx_incidents_active ON incidents(active);
CREATE INDEX idx_incidents_location ON incidents(latitude, longitude);
CREATE INDEX idx_location_checks_timestamp ON location_checks(timestamp);
CREATE INDEX idx_location_checks_user_id ON location_checks(user_id);
CREATE INDEX idx_location_checks_incident_id ON location_checks(incident_id);

-- Триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_incidents_updated_at 
    BEFORE UPDATE ON incidents 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Функция для поиска ближайших инцидентов
CREATE OR REPLACE FUNCTION find_nearby_incidents(
    p_latitude DECIMAL,
    p_longitude DECIMAL,
    p_radius_km DECIMAL
)
RETURNS TABLE(
    id BIGINT,
    user_id VARCHAR,
    latitude DECIMAL,
    longitude DECIMAL,
    title VARCHAR,
    description TEXT,
    severity VARCHAR,
    radius DECIMAL,
    active BOOLEAN,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    distance_km DECIMAL
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        i.*,
        (6371 * acos(
            cos(radians(p_latitude)) * cos(radians(i.latitude)) *
            cos(radians(i.longitude) - radians(p_longitude)) +
            sin(radians(p_latitude)) * sin(radians(i.latitude))
        )) as distance_km
    FROM incidents i
    WHERE i.active = true
    HAVING distance_km <= p_radius_km
    ORDER BY distance_km;
END;
$$ LANGUAGE plpgsql;