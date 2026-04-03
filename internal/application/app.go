package application

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"shb/internal/configs"
	"shb/internal/handlers"
	"shb/internal/repositories"
	"shb/internal/server"
	"shb/internal/services"
	"shb/pkg/db/cache/redisClient"
	"shb/pkg/db/pgx"
	"shb/pkg/external/email/smtpEmail"
	"shb/pkg/external/fs/minioFs"
	"shb/pkg/external/sms/smsProvider"
	"shb/pkg/logger"
	"shb/pkg/middlewares"
	"shb/pkg/rateLimiter/customLimiter"
	"shb/pkg/tokens/jwtToken"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

type App struct {
	config *configs.Config
	logger *logger.Logger
	server *server.Server
}

func NewApplication() *App {
	// 1. Load Internal Config
	cfg, err := configs.InitConfigs()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}
	// 2. Map Internal Config to Pkg Config for Logger
	// (Assuming pkgConfigs.Logger has a 'Level' field)
	log, err := logger.NewLogger(logger.Config{
		Level:         cfg.Logger.Level,
		Env:           cfg.App.Env,
		LogPath:       cfg.Logger.LogPath,
		IncludeCaller: cfg.Logger.IncludeCaller == "true",
	})
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

	// 3. SMS (Мапим конфиг)
	sms := smsProvider.NewSMSProvider(smsProvider.SMSConfig{
		APIKey:     cfg.SMS.APIKey,
		SenderName: cfg.SMS.SenderName,
		Login:      cfg.SMS.Login,
		BaseURL:    cfg.SMS.BaseURL,
	}, &log.Logger)

	// 4. Initialize SMTP Email Adapter
	emailAdapter := smtpEmail.NewSMTPEmail(&cfg.SMTP)

	token := jwtToken.NewJwtTokenIssuer(
		cfg.Security.JWTSecretKey,
		cfg.Security.AccessTokenTTL,
		cfg.Security.RefreshTokenTTL,
	)

	// 4. Middleware (Передаем секрет)
	middleware := middlewares.NewMiddleware(cfg.Security.JWTSecretKey)

	repository := repositories.NewRepository(postgresConn, &log.Logger)

	// Service uses Internal Config (ServiceConfig)
	// We pass emailAdapter here as it was required by Service constructor
	service := services.NewService(&cfg.Service, &log.Logger, repository, redis, sms, token, fileStorage, emailAdapter)

	handler := handlers.NewHandler(service, limiter, middleware, &log.Logger, cfg)

	// 5. Server (Мапим конфиг)
	readTimeout, _ := time.ParseDuration(cfg.Server.ReadTimeout)
	writeTimeout, _ := time.ParseDuration(cfg.Server.WriteTimeout)

	srv := server.NewServer(server.Config{
		Port:         cfg.Server.Port,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}, handler.InitRoutes()) // Внимание: NewServer в pkg/server теперь принимает handler http.Handler

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
		a.logger.Info().Msgf("%s is started on port %s", a.config.Server.Name, a.config.Server.Port)
		if err := a.server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Fatal().Msgf("Could not listen: %v", err)
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
