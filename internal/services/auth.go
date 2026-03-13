package services

import (
	"context"
	"errors"
	"fmt"
	"shb/internal/models"
	"shb/pkg/myerrors"
	"shb/pkg/utils"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) SendOTP(ctx context.Context, receiver string) (int, error) {
	otpCode, err := utils.GenerateOTP(s.cfg.Security.OTPLength)
	if err != nil {
		return 0, fmt.Errorf("otp code generation failed: %s", err)
	}

	expiresAt := time.Now().UTC().Add(s.cfg.Security.OTPDuration)
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

	go func(rcv, code, mthd string, id int) {
		var err error
		if mthd == "email" {
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
			err = s.email.SendEmail(context.Background(), rcv, subject, body)
		} else {
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

func (s *Service) ConfirmOTP(ctx context.Context, receiver, otp string) (*models.TokenResponse, error) {
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

	var user *models.User
	if strings.Contains(receiver, "@") {
		user, err = s.repo.GetUserByEmail(ctx, receiver)
	} else {
		user, err = s.repo.GetUserByPhone(ctx, receiver)
	}

	if err != nil {
		return nil, myerrors.NewBadRequestErr("Пользователь не найден")
	}

	if !user.IsActive {
		if err := s.repo.ActivateUser(ctx, user.ID); err != nil {
			return nil, err
		}
	}

	// Employee approval check (per user, not per institution)
	if user.Role == models.RoleEmployee && !user.IsApproved {
		return nil, myerrors.NewForbiddenErr("Ваш аккаунт ожидает подтверждения администратором")
	}

	access, refresh, err := s.token.IssueTokens(ctx, user.ID, user.Role, user.IsApproved)
	if err != nil {
		return nil, err
	}

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

func (s *Service) Login(ctx context.Context, email, password string) (*models.TokenResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, myerrors.ErrNotFound) {
			return nil, myerrors.NewUnauthorizedErr("Неверный пароль либо логин")
		}
		return nil, fmt.Errorf("get user error: %w", err)
	}

	if !user.IsActive {
		return nil, myerrors.NewUnauthorizedErr("Аккаунт не активирован. Пожалуйста, подтвердите OTP.")
	}

	// Employee approval check
	if user.Role == models.RoleEmployee && !user.IsApproved {
		return nil, myerrors.NewForbiddenErr("Ваш аккаунт ожидает подтверждения администратором")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password)); err != nil {
		return nil, myerrors.NewUnauthorizedErr("Неверный пароль либо логин")
	}

	access, refresh, err := s.token.IssueTokens(ctx, user.ID, user.Role, user.IsApproved)
	if err != nil {
		return nil, fmt.Errorf("issue tokens err: %w", err)
	}

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
	claims, err := s.token.VerifyToken(ctx, refreshToken)
	if err != nil {
		return nil, myerrors.NewUnauthorizedErr("Недействительный токен")
	}

	refreshHash := utils.HashToken(refreshToken)
	stored, err := s.repo.GetRefreshToken(ctx, refreshHash)
	if err != nil {
		return nil, myerrors.NewUnauthorizedErr("Токен не найден")
	}

	if stored.IsRevoked {
		_ = s.repo.RevokeAllUserRefreshTokens(ctx, stored.UserID)
		return nil, myerrors.NewForbiddenErr("Токен был отозван")
	}

	if time.Now().UTC().After(stored.ExpiresAt) {
		return nil, myerrors.NewUnauthorizedErr("Срок действия токена истек")
	}

	// Rotation
	if err := s.repo.RevokeRefreshToken(ctx, refreshHash); err != nil {
		return nil, fmt.Errorf("revoke old token: %w", err)
	}

	access, refresh, err := s.token.IssueTokens(ctx, claims.UserID, claims.Role, claims.IsApproved)
	if err != nil {
		return nil, fmt.Errorf("issue new tokens: %w", err)
	}

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
	email = strings.ToLower(strings.TrimSpace(email))
	existing, err := s.repo.GetUserByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil, myerrors.NewBadRequestErr("Пользователь с таким email уже существует")
	}

	if role == models.RoleEmployee {
		if institutionID == nil {
			return nil, myerrors.NewBadRequestErr("Для сотрудников обязательно указание учреждения")
		}
		inst, err := s.repo.GetInstitutionByID(ctx, *institutionID)
		if err != nil || inst.IsDeleted {
			return nil, myerrors.NewBadRequestErr("Указанное учреждение не найдено")
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	hashedStr := string(hashedPassword)

	newUser := &models.User{
		Email:         &email,
		Phone:         &phone,
		FullName:      &fullName,
		Password:      &hashedStr,
		Role:          role,
		InstitutionID: institutionID,
		IsActive:      false,
		IsApproved:    false, // Default to unapproved until superadmin approves
	}

	if err := s.repo.CreateUser(ctx, newUser); err != nil {
		return nil, err
	}

	if _, err := s.SendOTP(ctx, email); err != nil {
		s.logger.Error().Err(err).Msg("failed to send otp after register")
	}

	return nil, nil
}

func (s *Service) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return s.repo.GetUserByID(ctx, id)
}
