package services

import (
	"context"
	"errors"
	"fmt"
	"shb/internal/models"
	"shb/pkg/myerrors"
	"shb/pkg/utils"
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

	if err = s.repo.SaveOTP(ctx, otp); err != nil {
		return 0, err
	}

	// Асинхронная отправка
	go func(rcv, code, mthd string) {
		var err error
		if mthd == "email" {
			// Отправка Email
			err = s.email.SendOTP(rcv, code)
		} else {
			// Отправка SMS
			err = s.sms.SendSms(ctx, rcv, code)
		}

		if err != nil {
			s.logger.Error().Ctx(ctx).Err(err).
				Str("receiver", rcv).
				Str("method", mthd).
				Msg("failed to send otp")
		} else {
			s.logger.Info().Str("receiver", rcv).Str("method", mthd).Str("code", code).Msg("OTP sent (check logs)")
		}
	}(receiver, otpCode, method)

	return int(s.cfg.Security.OTPDuration.Seconds()), nil
}

// ConfirmOTP оставляем как был в твоем сообщении, он корректный
func (s *Service) ConfirmOTP(ctx context.Context, phone, otp string) (*models.TokenResponse, error) {
	otpDB, err := s.repo.GetOTP(ctx, phone)
	if err != nil {
		if errors.Is(err, myerrors.ErrNotFound) {
			return nil, myerrors.NewUnauthorizedErr("invalid OTP")
		}
		return nil, fmt.Errorf("get user by phone: %w", err)
	}

	if otp != otpDB.OTPCode {
		if err = s.repo.IncreaseOTPAttempt(ctx, otpDB.ID, phone); err != nil {
			s.logger.Error().Ctx(ctx).Err(err).Str("phone", phone).Msg("increase OTP attempt")
		}
		return nil, myerrors.NewUnauthorizedErr("invalid OTP")
	}

	if err = s.repo.MarkOTPAsVerified(ctx, otpDB.ID); err != nil {
		return nil, fmt.Errorf("mark otp as verified error: %w", err)
	}

	// ensure user exists
	user, err := s.ensureUserExists(ctx, phone)
	if err != nil {
		return nil, err
	}

	// issue tokens
	access, refresh, err := s.token.IssueTokens(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("issue tokens err: %w", err)
	}
	return &models.TokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

// ... методы Login, Register, ensureUserExists (оставляем старые) ...
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

	if err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password)); err != nil {
		return nil, myerrors.NewUnauthorizedErr("Неверный пароль либо логин")
	}

	access, refresh, err := s.token.IssueTokens(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("issue tokens err: %w", err)
	}
	return &models.TokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (s *Service) Register(ctx context.Context, email, phone, password, fullName, role string, institutionID *int) (*models.TokenResponse, error) {
	// 1. Проверяем email
	existing, err := s.repo.GetUserByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil, myerrors.NewBadRequestErr("Пользователь с таким email уже существует")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing failed: %w", err)
	}

	// ЛОГИКА ОБРАБОТКИ ТЕЛЕФОНА (Fix для unique constraint)
	var phonePtr *string
	if phone != "" {
		// Проверяем, не занят ли телефон, только если он указан
		existingPhone, err := s.repo.GetUserByPhone(ctx, phone)
		if err == nil && existingPhone != nil {
			return nil, myerrors.NewBadRequestErr("user with this phone already exists")
		}
		phonePtr = &phone
	} else {
		phonePtr = nil // Если пусто, шлем NULL в базу
	}

	hashedPasswordStr := string(hashedPassword)
	newUser := &models.User{
		Email:         &email,
		Phone:         phonePtr, // Используем указатель (nil или строка)
		Password:      &hashedPasswordStr,
		FullName:      &fullName,
		Role:          role,
		IsActive:      true,
		InstitutionID: institutionID,
	}

	if err := s.repo.CreateUser(ctx, newUser); err != nil {
		return nil, err
	}

	access, refresh, err := s.token.IssueTokens(ctx, newUser.ID)
	if err != nil {
		return nil, fmt.Errorf("issue tokens error: %w", err)
	}

	return &models.TokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
func (s *Service) GetUserByID(ctx context.Context, id int) (*models.User, error) {
    return s.repo.GetUserByID(ctx, id)
}