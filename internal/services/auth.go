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
	// 1. Проверки на существование (оставляем как есть)
	// ...

	// 2. Создаем пользователя НЕАКТИВНЫМ
	newUser := &models.User{
		// ... поля ...
		IsActive: false, // <--- ВАЖНО
	}

	if err := s.repo.CreateUser(ctx, newUser); err != nil {
		return nil, err
	}

	// 3. ОТПРАВЛЯЕМ КОД ПОДТВЕРЖДЕНИЯ
	// Используем email как receiver
	if _, err := s.SendOTP(ctx, email); err != nil {
		// Если не ушло — логируем, но юзера создали. 
		// (В идеале нужен роут "выслать код повторно")
		s.logger.Error().Err(err).Msg("failed to send otp after register")
	}

	// 4. Возвращаем nil, так как токенов еще нет
	return nil, nil
}
func (s *Service) GetUserByID(ctx context.Context, id int) (*models.User, error) {
    return s.repo.GetUserByID(ctx, id)
}