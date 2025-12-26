package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
    ServerPort string
    ServerHost string
    Environment string
    
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    DBSSLMode  string
    
    RedisHost     string
    RedisPort     string
    RedisPassword string
    RedisDB       int
    
    WebhookURL       string
    WebhookTimeout   time.Duration
    WebhookMaxRetries int
    WebhookRetryDelay time.Duration
    
    APIKeyOperator string
    
    StatsTimeWindowMinutes int
    CacheTTLMinutes       int
    LocationCheckRadiusKm float64
}

func Load() *Config {
    // Загрузка .env файла
    if err := godotenv.Load(); err != nil {
        log.Printf("No .env file found, using environment variables")
    }
    
    return &Config{
        ServerPort:  getEnv("SERVER_PORT", "8080"),
        ServerHost:  getEnv("SERVER_HOST", "0.0.0.0"),
        Environment: getEnv("ENVIRONMENT", "development"),
        
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", "postgres"),
        DBName:     getEnv("DB_NAME", "incident_system"),
        DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
        
        RedisHost:     getEnv("REDIS_HOST", "localhost"),
        RedisPort:     getEnv("REDIS_PORT", "6379"),
        RedisPassword: getEnv("REDIS_PASSWORD", ""),
        RedisDB:       getEnvAsInt("REDIS_DB", 0),
        
        WebhookURL:       getEnv("WEBHOOK_URL", "http://localhost:9090/webhook"),
        WebhookTimeout:   getEnvAsDuration("WEBHOOK_TIMEOUT", 5*time.Second),
        WebhookMaxRetries: getEnvAsInt("WEBHOOK_MAX_RETRIES", 3),
        WebhookRetryDelay: getEnvAsDuration("WEBHOOK_RETRY_DELAY", 1*time.Second),
        
        APIKeyOperator: getEnv("API_KEY_OPERATOR", "operator-key-secure-change-me"),
        
        StatsTimeWindowMinutes: getEnvAsInt("STATS_TIME_WINDOW_MINUTES", 60),
        CacheTTLMinutes:       getEnvAsInt("CACHE_TTL_MINUTES", 5),
        LocationCheckRadiusKm: getEnvAsFloat("LOCATION_CHECK_RADIUS_KM", 10.0),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    value := getEnv(key, "")
    if value == "" {
        return defaultValue
    }
    
    intValue, err := strconv.Atoi(value)
    if err != nil {
        log.Printf("Invalid value for %s: %v, using default", key, err)
        return defaultValue
    }
    
    return intValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
    value := getEnv(key, "")
    if value == "" {
        return defaultValue
    }
    
    floatValue, err := strconv.ParseFloat(value, 64)
    if err != nil {
        log.Printf("Invalid value for %s: %v, using default", key, err)
        return defaultValue
    }
    
    return floatValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
    value := getEnv(key, "")
    if value == "" {
        return defaultValue
    }
    
    duration, err := time.ParseDuration(value)
    if err != nil {
        log.Printf("Invalid duration for %s: %v, using default", key, err)
        return defaultValue
    }
    
    return duration
}