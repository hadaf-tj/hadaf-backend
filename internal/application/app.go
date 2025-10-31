package application

import (
	"context"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"os/signal"
	"shb/internal/handlers"
	"shb/internal/repositories"
	"shb/internal/services"
	"shb/pkg/configs"
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
)

type App struct {
	config *configs.Config
	logger *logger.Logger
	server *server.Server
}

func NewApplication() *App {
	cfg, err := configs.InitConfigs()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}
	log, err := logger.NewLogger(&cfg.Logger)
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
		Bucket:    configs.MinioBucket,
		Endpoint:  configs.MinioEndpoint,
		AccessKey: configs.MinioAccessKey,
		SecretKey: configs.MinioSecretKey,
		Logger:    &log.Logger,
	})
	if err != nil {
		panic("failed to initialize minio storage: " + err.Error())
	}

	limiter := customLimiter.NewRateLimiter(redis)
	sms := smsProvider.NewSMSProvider(&cfg.SMS)
	token := jwtToken.NewJwtTokenIssuer()
	middleware := middlewares.NewMiddleware()

	repository := repositories.NewRepository(postgresConn, &log.Logger)
	service := services.NewService(&cfg.Service, &log.Logger, repository, redis, sms, token, fileStorage)
	handler := handlers.NewHandler(service, limiter, middleware, &log.Logger, cfg)
	srv := server.NewServer(&cfg.Server, handler)

	return &App{
		config: cfg,
		logger: log,
		server: srv}
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
