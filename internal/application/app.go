package application

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"shb/internal/configs"
	"shb/internal/handlers"
	"shb/internal/repositories"
	"shb/internal/services"
	"shb/pkg/db/cache/redisClient"
	"shb/pkg/db/pgx"
	"shb/pkg/external/fs/minioFs"
	"shb/pkg/external/sms/smsProvider"
	"shb/pkg/logger"
	"shb/pkg/middlewares"
	"shb/pkg/rateLimiter/customLimiter"
	"shb/pkg/server"
	"shb/pkg/tokens/jwtToken"
	"syscall"
	"time"

	internalConfigs "shb/internal/configs" // Alias internal config

	"github.com/pkg/errors"
)

type App struct {
	config *configs.Config
	logger *logger.Logger
	server *server.Server
}

func NewApplication() *App {
	// 1. Load Internal Config
	cfg, err := internalConfigs.InitConfigs()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	// 2. Map Internal Config to Pkg Config for Logger
	// (Assuming pkgConfigs.Logger has a 'Level' field)
	legacyLoggerCfg := &internalConfigs.LoggerConfig{
		Level: cfg.Logger.Level,
	}
	log, err := logger.NewLogger(legacyLoggerCfg)
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}

	postgresConn, err := pgx.NewPgxPool()
	if err != nil {
		panic("failed to connect to postgres: " + err.Error())
	}

	redis, err := redisClient.NewRedisClient()
	if err != nil {
		panic("failed to connect to redis: " + err.Error())
	}

	fileStorage, err := minioFs.NewMinIOStorage(minioFs.MinIOConfig{
		Bucket:    cfg.Minio.Bucket,
		Endpoint:  cfg.Minio.Endpoint,
		AccessKey: cfg.Minio.AccessKey,
		SecretKey: cfg.Minio.SecretKey,
		Logger:    &log.Logger,
	})
	if err != nil {
		panic("failed to initialize minio storage: " + err.Error())
	}

	limiter := customLimiter.NewRateLimiter(redis)

	// 3. Map Internal Config to Pkg Config for SMS
	legacySMSCfg := &internalConfigs.SMSConfig{
		APIKey:     cfg.SMS.APIKey,
		SenderName: cfg.SMS.SenderName,
	}
	sms := smsProvider.NewSMSProvider(legacySMSCfg)

	// token := jwtToken.NewJwtTokenIssuer()
	token := jwtToken.NewJwtTokenIssuer(
		cfg.Security.JWTSecretKey,
		cfg.Security.AccessTokenTTL,
		cfg.Security.RefreshTokenTTL,
	)
	middleware := middlewares.NewMiddleware()

	repository := repositories.NewRepository(postgresConn, &log.Logger)

	// Service uses Internal Config (ServiceConfig)
	service := services.NewService(&cfg.Service, &log.Logger, repository, redis, sms, token, fileStorage)

	// Handler uses Internal Config
	handler := handlers.NewHandler(service, limiter, middleware, &log.Logger, cfg)

	// 4. Map Internal Config to Pkg Config for Server
	legacyServerCfg := &internalConfigs.ServerConfig{
		Name: cfg.Server.Name,
		Port: cfg.Server.Port,
	}
	srv := server.NewServer(legacyServerCfg, handler)

	return &App{
		config: cfg,
		logger: log,
		server: srv,
	}
}

func (a *App) Start() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		a.logger.Info().Msgf("%s is started", a.config.Server.Name)
		if err := a.server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Fatal().Msgf("Could not listen on %s: %v", a.config.Server.Port, err)
		}
	}()

	<-stopChan
	a.logger.Info().Msg("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Fatal().Msgf("Server forced to shutdown: %v", err)
	}

	a.logger.Info().Msg("Server exited properly")
}
