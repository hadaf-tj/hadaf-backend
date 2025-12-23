package configs

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	App      AppConfig
	Security SecurityConfig
	Database DatabaseConfig
	Logger   LoggerConfig
	SMS      SMSConfig
	Server   ServerConfig
	Service  ServiceConfig
}

type AppConfig struct {
	Port string
	Env  string
}

type SecurityConfig struct {
	JWTSecretKey            string
	AccessTokenTTL          time.Duration
	RefreshTokenTTL         time.Duration
	OTPLength               int
	OTPDuration             time.Duration
	OTPMaxAttempts          int
	OTPMaxAttemptsBlockTime time.Duration
	SendOTPAttempts         int
	SendOTPBlockTime        time.Duration
}

type DatabaseConfig struct {
	DSN string
}

type LoggerConfig struct {
	Level string
}

type SMSConfig struct {
	APIKey     string
	SenderName string
}

type ServerConfig struct {
	Name string
	Port string
}

type ServiceConfig struct {
	Security SecurityConfig
}

// Global constants (оставляем как есть, хотя лучше вынести в ENV)
const (
	MinioBucket    = "shb-files"
	MinioEndpoint  = "minio:9000"
	MinioAccessKey = "minioadmin"
	MinioSecretKey = "minioadmin"
)

// Helper для чтения ENV с дефолтным значением
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// InitConfigs loads the configuration
func InitConfigs() (*Config, error) {
	// Читаем переменные для PostgreSQL
	pgUser := getEnv("POSTGRES_USER", "postgres")
	pgPass := getEnv("POSTGRES_PASSWORD", "postgres")
	pgHost := getEnv("POSTGRES_HOST", "localhost") // В Docker будет "postgres"
	pgPort := getEnv("POSTGRES_PORT", "5432")
	pgDB := getEnv("POSTGRES_DB", "shb")

	// Формируем DSN строку
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		pgUser, pgPass, pgHost, pgPort, pgDB)

	// Читаем настройки Redis (важно для pkg/db/cache/redisClient)
	// Примечание: Если redisClient сам читает REDIS_HOST через os.Getenv, это сработает.
	// Если он берет конфиг отсюда - мы пока не передаем это явно в структуру,
	// но наличие переменных в ENV (через docker-compose) должно спасти ситуацию.

	return &Config{
		App: AppConfig{
			Port: getEnv("APP_PORT", ":8000"),
			Env:  getEnv("APP_ENV", "local"),
		},
		Security: SecurityConfig{
			JWTSecretKey:            getEnv("JWT_SECRET_KEY", "super_secret_dev_key"),
			AccessTokenTTL:          15 * time.Minute,
			RefreshTokenTTL:         720 * time.Hour,
			OTPLength:               4,
			OTPDuration:             5 * time.Minute,
			OTPMaxAttempts:          3,
			OTPMaxAttemptsBlockTime: 30 * time.Minute,
			SendOTPAttempts:         3,
			SendOTPBlockTime:        1 * time.Minute,
		},
		Database: DatabaseConfig{
			DSN: dsn, // Теперь DSN формируется динамически!
		},
		Logger: LoggerConfig{
			Level: getEnv("LOG_LEVEL", "debug"),
		},
		SMS: SMSConfig{
			APIKey:     getEnv("SMS_API_KEY", "mock"),
			SenderName: getEnv("SMS_SENDER_NAME", "Payvand"),
		},
		Server: ServerConfig{
			Name: "SocialHousingBackend",
			Port: getEnv("APP_PORT", ":8000"),
		},
		Service: ServiceConfig{
			Security: SecurityConfig{
				SendOTPAttempts:         3,
				SendOTPBlockTime:        1 * time.Minute,
				OTPMaxAttempts:          3,
				OTPMaxAttemptsBlockTime: 30 * time.Minute,
			},
		},
	}, nil
}