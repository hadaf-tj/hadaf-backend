package services

import (
	"context"
	"errors"
	"fmt"
	"shb/internal/models"
	"shb/pkg/myerrors"
	"shb/pkg/utils"
	"strconv"
	"strings" // Нужно для проверки @
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) SendOTP(ctx context.Context, receiver string) (int, error) {
	otpCode, err := utils.GenerateOTP(s.cfg.Security.OTPLength)
	if err != nil {
		return 0, fmt.Errorf("otp code generation failed: %s", err)
	}

	expiresAt := time.Now().UTC().Add(s.cfg.Security.OTPDuration)
	// ЛОГИКА ВЫБОРА КАНАЛА
	method := "sms"
	if strings.Contains(receiver, "@") {
		method = "email"
	}

	otp := &models.OTP{
		Receiver:   receiver,
		Method:     &method,
		OTPCode:    otpCode,
		SentAt:     time.Now().UTC(),
		ExpiresAt:  &expiresAt,
		IsVerified: false,
	}

	otpID, err := s.repo.SaveOTP(ctx, otp)
	if err != nil {
		return 0, err
	}

	// Асинхронная отправка
	go func(rcv, code, mthd string, id int) {
		var err error
		if mthd == "email" {
			// Отправка Email
			subject := "Ваш код подтверждения Hadaf"
			body := fmt.Sprintf(`
				<html>
				<body>
					<h2>Здравствуйте!</h2>
					<p>Ваш код подтверждения: <b>%s</b></p>
					<p>Никому не сообщайте этот код.</p>
				</body>
				</html>
			`, code)
			// Context is required but we are in goroutine. We can use Background or create new context.
			// Ideally pass ctx, but ctx might be cancelled. Use Background for async email.
			err = s.email.SendEmail(context.Background(), rcv, subject, body)
		} else {
			// Отправка SMS
			txnID := strconv.Itoa(id)
			err = s.sms.SendSms(context.Background(), rcv, code, txnID)
		}

		if err != nil {
			s.logger.Error().Ctx(ctx).Err(err).
				Str("receiver", rcv).
				Str("method", mthd).
				Msg("failed to send otp")
		} else {
			s.logger.Info().Str("receiver", rcv).Str("method", mthd).Msg("OTP sent successfully")
		}
	}(receiver, otpCode, method, otpID)

	return int(s.cfg.Security.OTPDuration.Seconds()), nil
}

// ConfirmOTP оставляем как был в твоем сообщении, он корректный
func (s *Service) ConfirmOTP(ctx context.Context, receiver, otp string) (*models.TokenResponse, error) {
	// 1. Проверяем OTP
	otpDB, err := s.repo.GetOTP(ctx, receiver)
	if err != nil {
		return nil, myerrors.NewUnauthorizedErr("Код не найден или истек")
	}

	if otp != otpDB.OTPCode {
		_ = s.repo.IncreaseOTPAttempt(ctx, otpDB.ID, receiver)
		return nil, myerrors.NewUnauthorizedErr("Неверный код")
	}

	if err = s.repo.MarkOTPAsVerified(ctx, otpDB.ID); err != nil {
		return nil, err
	}

	// 2. Ищем пользователя (по email или телефону)
	var user *models.User
	if strings.Contains(receiver, "@") {
		user, err = s.repo.GetUserByEmail(ctx, receiver)
	} else {
		user, err = s.repo.GetUserByPhone(ctx, receiver)
	}

	if err != nil {
		// Если пользователя нет — значит это просто проверка телефона/почты (например, для восстановления)
		// Но в нашем флоу регистрации пользователь уже создан (inactive).
		return nil, myerrors.NewBadRequestErr("Пользователь не найден")
	}

	// 3. АКТИВИРУЕМ ПОЛЬЗОВАТЕЛЯ
	if !user.IsActive {
		if err := s.repo.ActivateUser(ctx, user.ID); err != nil {
			return nil, err
		}
	}

	// 4. Выдаем токены
	access, refresh, err := s.token.IssueTokens(ctx, user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	// 5. Save refresh token to DB
	refreshHash := utils.HashToken(refresh)
	expiresAt := time.Now().UTC().Add(s.cfg.Security.RefreshTokenTTL)
	if err := s.repo.SaveRefreshToken(ctx, user.ID, refreshHash, expiresAt); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &models.TokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (s *Service) ensureUserExists(ctx context.Context, phone string) (*models.User, error) {
	user, err := s.repo.GetUserByPhone(ctx, phone)
	if err != nil && !errors.Is(err, myerrors.ErrNotFound) {
		return nil, fmt.Errorf("get user by phone: %w", err)
	}
	if user != nil {
		return user, nil
	}

	newUser := &models.User{Phone: &phone}
	err = s.repo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return newUser, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*models.TokenResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, myerrors.ErrNotFound) {
			return nil, myerrors.NewUnauthorizedErr("Неверный пароль либо логин")
		}
		return nil, fmt.Errorf("get user error: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password)); err != nil {
		return nil, myerrors.NewUnauthorizedErr("Неверный пароль либо логин")
	}

	access, refresh, err := s.token.IssueTokens(ctx, user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("issue tokens err: %w", err)
	}

	// Store refresh token in database (hashed)
	refreshHash := utils.HashToken(refresh)
	expiresAt := time.Now().UTC().Add(s.cfg.Security.RefreshTokenTTL)
	if err := s.repo.SaveRefreshToken(ctx, user.ID, refreshHash, expiresAt); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &models.TokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (s *Service) RefreshTokens(ctx context.Context, refreshToken string) (*models.TokenResponse, error) {
	// 1. Verify token signature and claims
	claims, err := s.token.VerifyToken(ctx, refreshToken)
	if err != nil {
		return nil, myerrors.NewUnauthorizedErr("Недействительный токен")
	}

	// 2. Hash and check if it exists and is not revoked
	refreshHash := utils.HashToken(refreshToken)
	stored, err := s.repo.GetRefreshToken(ctx, refreshHash)
	if err != nil {
		return nil, myerrors.NewUnauthorizedErr("Токен не найден")
	}

	if stored.IsRevoked {
		// Potential reuse attack!
		_ = s.repo.RevokeAllUserRefreshTokens(ctx, stored.UserID)
		return nil, myerrors.NewForbiddenErr("Токен был отозван")
	}

	if time.Now().UTC().After(stored.ExpiresAt) {
		return nil, myerrors.NewUnauthorizedErr("Срок действия токена истек")
	}

	// 3. Revoke current token (rotation)
	if err := s.repo.RevokeRefreshToken(ctx, refreshHash); err != nil {
		return nil, fmt.Errorf("revoke old token: %w", err)
	}

	// 4. Issue new pair
	access, refresh, err := s.token.IssueTokens(ctx, claims.UserID, claims.Role)
	if err != nil {
		return nil, fmt.Errorf("issue new tokens: %w", err)
	}

	// 5. Save new refresh token
	newRefreshHash := utils.HashToken(refresh)
	newExpiresAt := time.Now().UTC().Add(s.cfg.Security.RefreshTokenTTL)
	if err := s.repo.SaveRefreshToken(ctx, claims.UserID, newRefreshHash, newExpiresAt); err != nil {
		return nil, fmt.Errorf("save new refresh token: %w", err)
	}

	return &models.TokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (s *Service) RevokeAllUserRefreshTokens(ctx context.Context, userID int) error {
	return s.repo.RevokeAllUserRefreshTokens(ctx, userID)
}

func (s *Service) Register(ctx context.Context, email, phone, password, fullName, role string, institutionID *int) (*models.TokenResponse, error) {
	// 1. Проверяем, не существует ли уже пользователь с таким email
	email = strings.ToLower(strings.TrimSpace(email))
	existing, err := s.repo.GetUserByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil, myerrors.NewBadRequestErr("Пользователь с таким email уже существует")
	}

	// 2. Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	hashedStr := string(hashedPassword)

	// 3. Создаем пользователя НЕАКТИВНЫМ (со всеми полями)
	newUser := &models.User{
		Email:         &email,
		Phone:         &phone,
		FullName:      &fullName,
		Password:      &hashedStr,
		Role:          role,
		InstitutionID: institutionID,
		IsActive:      false,
	}

	if err := s.repo.CreateUser(ctx, newUser); err != nil {
		return nil, err
	}

	// 4. ОТПРАВЛЯЕМ КОД ПОДТВЕРЖДЕНИЯ
	if _, err := s.SendOTP(ctx, email); err != nil {
		s.logger.Error().Err(err).Msg("failed to send otp after register")
	}

	// 5. Возвращаем nil, так как токенов еще нет (нужно подтвердить OTP)
	return nil, nil
}

func (s *Service) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return s.repo.GetUserByID(ctx, id)
}
