package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"incident-system/internal/config"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
    db *sql.DB
}

func NewPostgresDB(cfg *config.Config) (*PostgresDB, error) {
    connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
    )
    
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    
    // Настройка пула соединений
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    // Проверка подключения
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := db.PingContext(ctx); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    log.Println("Successfully connected to PostgreSQL")
    return &PostgresDB{db: db}, nil
}

func (p *PostgresDB) Close() error {
    return p.db.Close()
}

func (p *PostgresDB) GetDB() *sql.DB {
    return p.db
}