package services

import (
	"context"
	"errors"
	"fmt"
	"shb/internal/models"
	"shb/pkg/myerrors"
	"shb/pkg/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) SendOTP(ctx context.Context, receiver string) (int, error) {
	otpCode, err := utils.GenerateOTP(s.cfg.Security.OTPLength)
	if err != nil {
		return 0, fmt.Errorf("otp code generation failed: %s", err)
	}

	method := "sms"
	expiresAt := time.Now().UTC().Add(s.cfg.Security.OTPDuration)
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

	go func(rcv, code string) {
		if err = s.sms.SendSms(ctx, rcv, code); err != nil {
			s.logger.Error().Ctx(ctx).Err(err).Str("receiver", rcv).Msg("send sms failed")
		}
	}(receiver, otpCode)

	return int(s.cfg.Security.OTPDuration.Seconds()), nil
}

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

func (s *Service) Login(ctx context.Context, phone, password string) (*models.TokenResponse, error) {
	user, err := s.repo.GetUserByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, myerrors.ErrNotFound) {
			return nil, myerrors.NewUnauthorizedErr("invalid phone_number or password")
		}
		return nil, fmt.Errorf("get user by phone: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password)); err != nil {
		//todo нужно ли блокировать при большом количестве неудачных попытках
		return nil, myerrors.NewUnauthorizedErr("invalid phone_number or password")
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

func (s *Service) Register(ctx context.Context, phone, password, fullName string, institutionID int) (*models.TokenResponse, error) {
	// 1. Проверяем, есть ли уже такой пользователь
	existing, err := s.repo.GetUserByPhone(ctx, phone)
	if err == nil && existing != nil {
		return nil, myerrors.NewBadRequestErr("user with this phone already exists")
	}

	// 2. Хешируем пароль (Правильно, средствами Go!)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing failed: %w", err)
	}

	// 3. Создаем модель
	hashedPasswordStr := string(hashedPassword)
	newUser := &models.User{
		Phone:         &phone,
		Password:      &hashedPasswordStr,
		FullName:      &fullName,
		Role:          models.RoleEmployee, // По умолчанию создаем Сотрудника
		IsActive:      true,
		InstitutionID: &institutionID,
	}

	// 4. Сохраняем в базу
	if err := s.repo.CreateUser(ctx, newUser); err != nil {
		return nil, err
	}

	// 5. Сразу логиним его (выдаем токены)
	access, refresh, err := s.token.IssueTokens(ctx, newUser.ID)
	if err != nil {
		return nil, fmt.Errorf("issue tokens error: %w", err)
	}

	return &models.TokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
