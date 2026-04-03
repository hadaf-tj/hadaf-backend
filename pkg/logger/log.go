package logger

import (
	"os"
	"shb/pkg/constants"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
}

// Config определяем прямо тут, чтобы не зависеть от внешних пакетов
type Config struct {
	Level         string
	Env           string // "local", "prod", etc.
	LogPath       string // путь к файлу, напр. "logs/app.log"
	IncludeCaller bool
}

var instance *Logger

func NewLogger(cfg Config) (*Logger, error) {
	var output *os.File

	if cfg.Env == constants.LocalAppEnv && cfg.LogPath != "" {
		f, err := os.OpenFile(cfg.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			return nil, err
		}
		output = f
	} else {
		output = os.Stdout
	}

	zerolog.TimeFieldFormat = time.RFC3339

	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.DebugLevel
	}

	zl := zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Logger()

	if cfg.IncludeCaller {
		zl = zl.With().Caller().Logger()
	}

	instance = &Logger{zl}
	return instance, nil
}
