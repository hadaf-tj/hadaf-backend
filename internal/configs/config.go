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
	SMTP     SMTPConfig
	Server   ServerConfig
	Service  ServiceConfig
	Redis    RedisConfig
	Minio    MinioConfig
}

type AppConfig struct {
	Port string
	Env  string
}

type SecurityConfig struct {
	JWTSecretKey            string
	AccessTokenTTL          time.Duration
	AccessTokenSecret       string
	RefreshTokenTTL         time.Duration
	RefreshTokenSecret      string
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
	Level         string
	LogPath       string
	IncludeCaller string
}

type SMSConfig struct {
	APIKey     string
	SenderName string
	Login      string
	BaseURL    string
}

type SMTPConfig struct {
	Host      string
	Port      string
	Username  string
	Password  string
	FromEmail string
	FromName  string
}

type ServerConfig struct {
	Name         string
	Port         string
	Host         string
	WriteTimeout string
	ReadTimeout  string
}

type ServiceConfig struct {
	Security SecurityConfig
}

type RedisConfig struct {
	Host      string
	Port      string
	DefaultDB int
	Timeout   string
}

type MinioConfig struct {
	Bucket    string
	Endpoint  string
	AccessKey string
	SecretKey string
}

// Helper для чтения ENV с дефолтным значением
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// requireEnv panics if env variable is missing or empty (for secrets)
func requireEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		panic(fmt.Sprintf("FATAL: required env var %s is not set", key))
	}
	return v
}

// InitConfigs loads the configuration
func InitConfigs() (*Config, error) {
	// Читаем переменные для PostgreSQL
	pgUser := getEnv("POSTGRES_USER", "postgres")
	pgPass := requireEnv("POSTGRES_PASSWORD")
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
			JWTSecretKey:            requireEnv("JWT_SECRET_KEY"),
			AccessTokenTTL:          15 * time.Minute,
			AccessTokenSecret:       requireEnv("ACCESS_TOKEN_SECRET"),
			RefreshTokenTTL:         720 * time.Hour,
			RefreshTokenSecret:      requireEnv("REFRESH_TOKEN_SECRET"),
			OTPLength:               6,
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
			Level:         getEnv("LOG_LEVEL", "debug"),
			LogPath:       getEnv("LOG_PATH", ""),
			IncludeCaller: getEnv("INCLUDE_CALLER", ""),
		},
		SMS: SMSConfig{
			APIKey:     getEnv("SMS_API_KEY", "mock"),
			SenderName: getEnv("SMS_SENDER_NAME", "Payvand"),
			Login:      getEnv("SMS_LOGIN", ""),
			BaseURL:    getEnv("SMS_BASE_URL", "https://api.osonsms.com"),
		},
		SMTP: SMTPConfig{
			Host:      getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port:      getEnv("SMTP_PORT", "587"),
			Username:  getEnv("SMTP_USERNAME", ""),
			Password:  getEnv("SMTP_PASSWORD", ""),
			FromEmail: getEnv("SMTP_FROM_EMAIL", "noreply@socialhousing.tj"),
			FromName:  getEnv("SMTP_FROM_NAME", "Social Housing Platform"),
		},
		Server: ServerConfig{
			Name:         "SocialHousingBackend",
			Port:         getEnv("APP_PORT", ":8000"),
			Host:         getEnv("APP_HOST", "localhost"),
			WriteTimeout: getEnv("APP_WRITE_TIMEOUT", "10s"),
			ReadTimeout:  getEnv("APP_READ_TIMEOUT", "10s"),
		},
		Service: ServiceConfig{
			Security: SecurityConfig{
				JWTSecretKey:            requireEnv("JWT_SECRET_KEY"),
				OTPLength:               6,
				OTPDuration:             5 * time.Minute,
				OTPMaxAttempts:          3,
				OTPMaxAttemptsBlockTime: 30 * time.Minute,
				SendOTPAttempts:         3,
				SendOTPBlockTime:        1 * time.Minute,
			},
		},
		Redis: RedisConfig{
			Host:      getEnv("REDIS_HOST", "localhost"),
			Port:      getEnv("REDIS_PORT", "6379"),
			DefaultDB: 0,
			Timeout:   getEnv("REDIS_TIMEOUT", "5s"),
		},
		Minio: MinioConfig{
			Bucket:    getEnv("MINIO_BUCKET", "minio"),
			Endpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", "minio"),
			SecretKey: getEnv("MINIO_SECRET_KEY", "minio"),
		},
	}, nil
}
