package services

import (
	"context"
	"shb/internal/configs"
	"shb/internal/models"
	"shb/internal/repositories/filters"
	"shb/pkg/db/cache"
	"shb/pkg/external/email"
	"shb/pkg/external/fs"
	"shb/pkg/external/sms"
	"shb/pkg/tokens"

	"github.com/rs/zerolog"
)

// IRepository описывает методы доступа к БД.
type IRepository interface {
	// GetUserByPhone возвращает пользователя по номеру телефона.
	GetUserByPhone(ctx context.Context, phone string) (*models.User, error)
	// CreateUser создаёт нового пользователя.
	CreateUser(ctx context.Context, user *models.User) error

	// SaveOTP сохраняет новый OTP-код в базу данных.
	SaveOTP(ctx context.Context, o *models.OTP) error
	// GetOTP получает последний активный и неподтверждённый OTP-код по номеру телефона.
	GetOTP(ctx context.Context, phone string) (*models.OTP, error)
	// MarkOTPAsVerified отмечает OTP-код как подтверждённый.
	MarkOTPAsVerified(ctx context.Context, otpID int) error
	// IncreaseOTPAttempt увеличивает счётчик попыток ввода OTP-кода.
	IncreaseOTPAttempt(ctx context.Context, otpID int, phone string) error

	// --- Institution Methods ---
	GetAllInstitutions(ctx context.Context, filter filters.InstitutionFilter) ([]*models.Institution, error)
	CreateInstitution(ctx context.Context, i *models.Institution) (int, error)
	GetInstitutionByID(ctx context.Context, id int) (*models.Institution, error)

	// --- Needs Methods ---
	CreateNeed(ctx context.Context, need *models.Need) (int, error)
	GetNeedByID(ctx context.Context, id int) (*models.Need, error)
	UpdateNeed(ctx context.Context, n *models.Need) error
	DeleteNeed(ctx context.Context, id int) error
	GetNeedsByInstitution(ctx context.Context, filter filters.NeedsFilter, institutionID int) ([]*models.Need, error)

	GetUserByID(ctx context.Context, id int) (*models.User, error)

	// --- Booking Methods ---
	CreateBooking(ctx context.Context, booking *models.Booking) (int, error)
	GetBookingByID(ctx context.Context, id int) (*models.Booking, error)
	GetBookingsByNeed(ctx context.Context, needID int) ([]*models.Booking, error)
	GetBookingsByUser(ctx context.Context, userID int) ([]*models.Booking, error)
	UpdateBookingStatus(ctx context.Context, bookingID int, status string) error
	GetBookingsByInstitution(ctx context.Context, institutionID int) ([]*models.Booking, error)

	// --- Event Methods ---
	CreateEvent(ctx context.Context, e *models.Event) (int, error)
	GetEventByID(ctx context.Context, id int) (*models.Event, error)
	GetAllEvents(ctx context.Context, userID int) ([]*models.EventResponse, error)
	JoinEvent(ctx context.Context, eventID, userID int) error
	LeaveEvent(ctx context.Context, eventID, userID int) error
}
type Service struct {
	cfg    *configs.ServiceConfig // CHANGED: from configs.Service to configs.ServiceConfig
	logger *zerolog.Logger
	repo   IRepository
	cache  cache.ICache
	sms    sms.ISmsAdapter
	token  tokens.ITokenIssuer
	fs     fs.Storage
	email  email.IEmailAdapter
}

func NewService(cfg *configs.ServiceConfig, log *zerolog.Logger, repo IRepository, cache cache.ICache,
	sms sms.ISmsAdapter, token tokens.ITokenIssuer, fs fs.Storage, email email.IEmailAdapter) *Service {
	return &Service{
		cfg:    cfg,
		logger: log,
		repo:   repo,
		cache:  cache,
		sms:    sms,
		token:  token,
		fs:     fs,
		email:  email}
}
