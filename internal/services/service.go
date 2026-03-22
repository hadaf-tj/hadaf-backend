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
	"time"

	"github.com/rs/zerolog"
)

// IRepository описывает методы доступа к БД.
type IRepository interface {
	// GetUserByPhone возвращает пользователя по номеру телефона.
	GetUserByPhone(ctx context.Context, phone string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	// CreateUser создаёт нового пользователя.
	CreateUser(ctx context.Context, user *models.User) error
	ActivateUser(ctx context.Context, id int) error

	// SaveOTP сохраняет новый OTP-код в базу данных.
	SaveOTP(ctx context.Context, o *models.OTP) (int, error)
	// GetOTP получает последний активный и неподтверждённый OTP-код по номеру телефона.
	GetOTP(ctx context.Context, phone string) (*models.OTP, error)
	// MarkOTPAsVerified отмечает OTP-код как подтверждённый.
	MarkOTPAsVerified(ctx context.Context, otpID int) error
	// IncreaseOTPAttempt увеличивает счётчик попыток ввода OTP-кода.
	IncreaseOTPAttempt(ctx context.Context, otpID int, phone string) error

	// --- Institution Methods ---
	GetAllInstitutions(ctx context.Context, search string, iType string, userLat, userLng float64, sortBy string) ([]*models.Institution, error)
	CreateInstitution(ctx context.Context, i *models.Institution) (int, error)
	GetInstitutionByID(ctx context.Context, id int) (*models.Institution, error)

	// --- Needs Methods ---
	CreateNeed(ctx context.Context, need *models.Need) (int, error)
	GetNeedByID(ctx context.Context, id int) (*models.Need, error)
	UpdateNeed(ctx context.Context, n *models.Need) error
	DeleteNeed(ctx context.Context, id int) error
	GetNeedsByInstitution(ctx context.Context, filter filters.NeedsFilter, institutionID int) ([]*models.Need, error)

	// --- Booking Methods ---
	CreateBooking(ctx context.Context, booking *models.Booking) (int, error)
	GetBookingByID(ctx context.Context, id int) (*models.Booking, error)
	GetBookingsByNeed(ctx context.Context, needID int) ([]*models.Booking, error)
	GetBookingsByUser(ctx context.Context, userID int) ([]*models.Booking, error)
	GetActiveBookingByUserAndNeed(ctx context.Context, userID, needID int) (*models.Booking, error)
	UpdateBookingStatus(ctx context.Context, bookingID int, status string) error
	UpdateBookingQuantity(ctx context.Context, bookingID int, qty float64) error
	IncrementReceivedQty(ctx context.Context, needID int, qty float64) error
	GetBookingsByInstitution(ctx context.Context, institutionID int) ([]*models.Booking, error)

	// --- Event Methods ---
	CreateEvent(ctx context.Context, e *models.Event) (int, error)
	GetEventByID(ctx context.Context, id int) (*models.Event, error)
	GetAllEvents(ctx context.Context, userID int) ([]*models.EventResponse, error)
	JoinEvent(ctx context.Context, eventID, userID int) error
	LeaveEvent(ctx context.Context, eventID, userID int) error
	GetInstitutionEvents(ctx context.Context, institutionID int) ([]*models.EventResponse, error)
	UpdateEventStatus(ctx context.Context, eventID int, status string) error

	// Vacancies
	GetAllVacancies(ctx context.Context) ([]*models.Vacancy, error)
	GetVacancyByID(ctx context.Context, id int) (*models.Vacancy, error)

	// Team Members
	GetAllTeamMembers(ctx context.Context) ([]*models.TeamMember, error)
	GetTeamMemberByID(ctx context.Context, id int) (*models.TeamMember, error)

	// --- Stats Methods ---
	GetPublicStats(ctx context.Context) (map[string]int, error)

	CreateNeedHistory(ctx context.Context, history *models.NeedsHistory) error

	// --- Token Methods ---
	SaveRefreshToken(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, tokenHash string) (*models.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
	RevokeAllUserRefreshTokens(ctx context.Context, userID int) error
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
		email:  email,
		token:  token,
		fs:     fs,
	}
}
