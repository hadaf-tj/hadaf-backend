package configs

import "time"

type Config struct {
	App      AppConfig
	Security SecurityConfig
	Database DatabaseConfig
	Logger   LoggerConfig   // Added
	SMS      SMSConfig      // Added
	Server   ServerConfig   // Added
	Service  ServiceConfig  // Added
}

type AppConfig struct {
	Port string
	Env  string
}

type SecurityConfig struct {
	JWTSecretKey            string        // Fixed name (was SecretKey)
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

// Global constants
const (
	MinioBucket    = "shb-files"
	MinioEndpoint  = "minio:9000"
	MinioAccessKey = "minioadmin"
	MinioSecretKey = "minioadmin"
)

// InitConfigs loads the configuration
func InitConfigs() (*Config, error) {
	// Hardcoded for now to pass compilation.
	// In real app, load from .env using Viper or Godotenv
	return &Config{
		App: AppConfig{Port: ":8000", Env: "local"},
		Security: SecurityConfig{
			JWTSecretKey:            "super_secret_key",
			AccessTokenTTL:          15 * time.Minute,
			RefreshTokenTTL:         720 * time.Hour,
			OTPLength:               4,
			OTPDuration:             5 * time.Minute,
			OTPMaxAttempts:          3,
			OTPMaxAttemptsBlockTime: 30 * time.Minute,
			SendOTPAttempts:         3,
			SendOTPBlockTime:        1 * time.Minute,
		},
		Database: DatabaseConfig{DSN: ""},
		Logger:   LoggerConfig{Level: "debug"},
		SMS:      SMSConfig{APIKey: "mock", SenderName: "Test"},
		Server:   ServerConfig{Name: "SocialHousingBackend", Port: ":8000"},
		Service: ServiceConfig{
			Security: SecurityConfig{
				SendOTPAttempts:  3,
				SendOTPBlockTime: 1 * time.Minute,
                OTPMaxAttempts: 3,
                OTPMaxAttemptsBlockTime: 30 * time.Minute,
			},
		},
	}, nil
}