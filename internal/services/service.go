package services

import (
	"context"
	"shb/internal/models"
	"shb/pkg/configs"
	"shb/pkg/db/cache"
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
	// GetAllInstitutions возвращает список учреждений с возможностью фильтрации по городу.
	GetAllInstitutions(ctx context.Context, city string) ([]*models.Institution, error)
	// CreateInstitution создает новое учреждение
	CreateInstitution(ctx context.Context, i *models.Institution) (int, error)
}

type Service struct {
	cfg    *configs.Service
	logger *zerolog.Logger
	repo   IRepository
	cache  cache.ICache
	sms    sms.ISmsAdapter
	token  tokens.ITokenIssuer
	fs     fs.Storage
}

func NewService(cfg *configs.Service, log *zerolog.Logger, repo IRepository, cache cache.ICache,
	sms sms.ISmsAdapter, token tokens.ITokenIssuer, fs fs.Storage) *Service {
	return &Service{
		cfg:    cfg,
		logger: log,
		repo:   repo,
		cache:  cache,
		sms:    sms,
		token:  token,
		fs:     fs}
}