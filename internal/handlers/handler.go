package handlers

import (
	"context"
	"errors"
	"net/http"
	"shb/internal/configs"
	"shb/internal/models"
	"shb/internal/repositories/filters"
	"shb/pkg/constants"
	"shb/pkg/middlewares"
	"shb/pkg/myerrors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Limiter interface {
	Allow(ctx context.Context, key string, limit int, windowSeconds int) (bool, error)
	ResetAttempts(ctx context.Context, key string) error
}
type IService interface {
	SendOTP(ctx context.Context, receiver string) (int, error)
	ConfirmOTP(ctx context.Context, phone, otp string) (*models.TokenResponse, error)
	Login(ctx context.Context, phone, password string) (*models.TokenResponse, error)
	Register(ctx context.Context, email, phone, password, fullName, role string, institutionID *int) (*models.TokenResponse, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)

	GetAllInstitutions(ctx context.Context, search string, iType string, userLat, userLng float64, sortBy string) ([]*models.Institution, error)
	CreateInstitution(ctx context.Context, i *models.Institution) (int, error)
	GetInstitutionByID(ctx context.Context, id int) (*models.Institution, error)
	

	CreateNeed(ctx context.Context, need *models.Need) (int, error)
	UpdateNeed(ctx context.Context, n *models.Need) error
	DeleteNeed(ctx context.Context, id int) error
	GetNeedsByInstitution(ctx context.Context, filter filters.NeedsFilter, institutionID int) ([]*models.Need, error)
}

type Handler struct {
	service    IService
	limiter    Limiter                 // CHANGED: Use local interface
	middleware *middlewares.Middleware // CHANGED: Use imported type (pointer likely)
	logger     *zerolog.Logger
	cfg        *configs.Config
}

func NewHandler(service IService, limiter Limiter, middleware *middlewares.Middleware, logger *zerolog.Logger, cfg *configs.Config) *Handler {
	return &Handler{
		service:    service,
		limiter:    limiter,
		middleware: middleware,
		logger:     logger,
		cfg:        cfg,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(h.CORSMiddleware(), gin.RecoveryWithWriter(gin.DefaultWriter), h.RequestID())
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
		v1.POST("/register", h.register)

		// Исправленный вызов middleware
		v1.GET("/check_access", h.middleware.AuthMiddleware(), func(c *gin.Context) {
			h.success(c, "valid")
		})
		v1.GET("/me", h.middleware.AuthMiddleware(), h.getMe)

		v1.GET("/institutions", h.getAllInstitutions)
		v1.GET("/institutions/:id", h.getInstitutionByID)
		v1.POST("/institutions", h.createInstitution)

		v1.GET("/institutions/:id/needs", h.getNeedsByInstitution)

		needs := v1.Group("/needs")
		needs.Use(h.middleware.AuthMiddleware(models.RoleEmployee, models.RoleSuperAdmin))
		{
			needs.POST("", h.createNeed)
			needs.PUT("/:id", h.updateNeed)
			needs.DELETE("/:id", h.deleteNeed)
		}
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

func (h *Handler) RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get(constants.RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, constants.RequestIDKey, requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func (h *Handler) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, x-request-id")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}