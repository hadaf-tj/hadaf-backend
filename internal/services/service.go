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
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	// CreateUser создаёт нового пользователя.
	CreateUser(ctx context.Context, user *models.User) error
	ActivateUser(ctx context.Context, id int) error
	

	// SaveOTP сохраняет новый OTP-код в базу данных.
	SaveOTP(ctx context.Context, o *models.OTP) error
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
	CreateNeedHistory(ctx context.Context, history *models.NeedsHistory) error
}
type Service struct {
	cfg    *configs.ServiceConfig // CHANGED: from configs.Service to configs.ServiceConfig
	logger *zerolog.Logger
	repo   IRepository
	cache  cache.ICache
	sms    sms.ISmsAdapter
	token  tokens.ITokenIssuer
	fs     fs.Storage
	email  email.ISender
}

func NewService(cfg *configs.ServiceConfig, log *zerolog.Logger, repo IRepository, cache cache.ICache,
    sms sms.ISmsAdapter, token tokens.ITokenIssuer, fs fs.Storage) *Service {
		emailSender := email.NewSmtpSender()
	return &Service{
		cfg:    cfg,
		logger: log,
		repo:   repo,
		cache:  cache,
		sms:    sms,
		email:  emailSender,
		token:  token,
		fs:     fs}
}
