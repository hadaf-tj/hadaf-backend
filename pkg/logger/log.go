package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"shb/pkg/configs"
	"shb/pkg/constants"
	"strings"
	"time"
)

type Logger struct {
	zerolog.Logger
}

var instance *Logger

func NewLogger(cfg *configs.Logger) (*Logger, error) {
	var (
		output *os.File
		err    error
	)

	if cfg.LogPath != "" {
		output, err = os.OpenFile(cfg.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
	} else {
		output = os.Stdout
	}

	level, err := zerolog.ParseLevel(strings.ToLower(cfg.Level))
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}
	zerolog.SetGlobalLevel(level)

	zerolog.TimeFieldFormat = time.RFC3339

	zl := zerolog.New(output).
		Hook(RequestIDHook{}).
		With().
		Timestamp().
		Logger()

	if cfg.IncludeCaller {
		zl = zl.With().Caller().Logger()
	}

	instance = &Logger{Logger: zl}

	return instance, nil
}

type RequestIDHook struct{}

func (RequestIDHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	ctx := e.GetCtx()
	if ctx == nil {
		return
	}
	if requestID, ok := ctx.Value(constants.RequestIDKey).(string); ok && requestID != "" {
		e.Str("request_id", requestID)
	}
}
