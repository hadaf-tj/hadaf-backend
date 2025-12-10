package configs

import (
	"time"

	"github.com/spf13/viper"
)

func InitConfigs() (*Config, error) {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err = viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

type Config struct {
	App     App
	Security SecurityConfig
	Server  Server
	Service Service
	Logger  Logger
	SMS     SMSProvider
}

type App struct {
}
type AppConfig struct {
    Port string
}

type SecurityConfig struct {
    SecretKey string
}

type DatabaseConfig struct {
    DSN string
}

type Server struct {
	Name         string
	Host         string
	Port         string
	WriteTimeout int64
	ReadTimeout  int64
}

type Logger struct {
	LogPath       string
	Level         string
	IncludeCaller bool
}

type Service struct {
	Security Security
}

type Security struct {
	SendOTPAttempts         int
	SendOTPBlockTime        int
	OTPLength               int
	OTPMaxAttempts          int
	MaxLoginAttempts        int
	OTPMaxAttemptsBlockTime time.Duration
	OTPDuration             time.Duration
}

type SMSProvider struct {
    APIKey     string // <--- Этого поля не хватало
    SenderName string // <--- И этого тоже
}
