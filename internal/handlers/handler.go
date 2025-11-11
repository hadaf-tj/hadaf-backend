package handlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"shb/internal/models"
	"shb/pkg/configs"
	"shb/pkg/middlewares"
	"shb/pkg/myerrors"
	"shb/pkg/rateLimiter"
)

// IService описывает бизнес-логику.
type IService interface {
	// SendOTP отправляет OTP на указанный номер телефона.
	SendOTP(ctx context.Context, receiver string) (int, error)
	// ConfirmOTP проверяет OTP и выдаёт токен при успешной верификации.
	ConfirmOTP(ctx context.Context, phone, otp string) (*models.TokenResponse, error)
	// Login проверяет логин и пароль, выдаёт токен при успешной верификации.
	Login(ctx context.Context, phone, password string) (*models.TokenResponse, error)
}

type Handler struct {
	service    IService
	limiter    rateLimiter.IRateLimiter
	middleware *middlewares.Middleware
	cfg        *configs.Config
	logger     *zerolog.Logger
}

func NewHandler(service IService, limiter rateLimiter.IRateLimiter,
	middleware *middlewares.Middleware, log *zerolog.Logger, cfg *configs.Config) *Handler {
	return &Handler{
		service:    service,
		limiter:    limiter,
		middleware: middleware,
		cfg:        cfg,
		logger:     log,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(h.middleware.CORSMiddleware(), gin.RecoveryWithWriter(gin.DefaultWriter), h.RequestID())
	router.NoRoute(h.noRoute)
	router.GET("/ping", h.ping)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
			ginSwagger.URL("/api/v1/docs/swagger.yaml"),
		))
		v1.Static("/docs", "./docs")

		v1.POST("/send_otp", h.sendOTP)
		v1.POST("/confirm_otp", h.confirmOTP)
		v1.POST("/login", h.login)

		v1.GET("/check_access", h.middleware.AccessToken())
		v1.GET("/check_refresh", h.middleware.RefreshToken())
	}
	return router
}

func (h *Handler) ping(context *gin.Context) {
	h.respond(context, "pong", http.StatusOK)
}

func (h *Handler) noRoute(context *gin.Context) {
	h.respond(context, "this route is not supported", http.StatusNotFound)
}

func (h *Handler) respond(context *gin.Context, obj interface{}, code int) {
	context.JSON(code, obj)
}

func (h *Handler) success(c *gin.Context, data any) {
	h.respond(c, models.Response{
		Message: "Success",
		Data:    data,
	}, http.StatusOK)
}

func (h *Handler) handleError(c *gin.Context, err error) {
	badReq := &myerrors.BadRequestErr{}
	forbidden := &myerrors.ForbiddenErr{}
	unprocessable := &myerrors.UnprocessableErr{}
	unauth := &myerrors.UnauthorizedErr{}
	manyReq := &myerrors.TooManyRequestsErr{}

	switch {
	case errors.As(err, unprocessable):
		c.JSON(http.StatusUnprocessableEntity, unprocessable)
	case errors.As(err, badReq):
		c.JSON(http.StatusBadRequest, badReq)
	case errors.As(err, forbidden):
		c.JSON(http.StatusForbidden, forbidden)
	case errors.As(err, unauth):
		c.JSON(http.StatusUnauthorized, unauth)
	case errors.As(err, manyReq):
		c.JSON(http.StatusTooManyRequests, manyReq)
	default:
		c.JSON(http.StatusInternalServerError, myerrors.InternalError())
	}
	c.Abort()
}
