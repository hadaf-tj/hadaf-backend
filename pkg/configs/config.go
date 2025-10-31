package configs

import (
	"github.com/spf13/viper"
	"time"
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
	Server  Server
	Service Service
	Logger  Logger
	SMS     SMSProvider
}

type App struct {
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
type SMSProvider struct{}
