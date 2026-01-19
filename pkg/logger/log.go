package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
}

// Config определяем прямо тут, чтобы не зависеть от внешних пакетов
type Config struct {
	Level         string
	IncludeCaller bool
}

var instance *Logger

func NewLogger(cfg Config) (*Logger, error) {
	var (
		output *os.File
		err    error
	)

	// ... логика открытия файла логов, если нужно ...
	output = os.Stdout

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