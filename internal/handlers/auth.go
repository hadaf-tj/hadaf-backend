package handlers

import (
	"fmt"
	"shb/pkg/myerrors"
	"shb/pkg/utils"

	"github.com/gin-gonic/gin"
)

func (h *Handler) sendOTP(c *gin.Context) {
	ctx := c.Request.Context()
	logger := h.logger.With().
		Ctx(ctx).
		Str("handler", "sendOTP").
		Logger()

	in := struct {
		Receiver     string `json:"receiver" binding:"required"`
		CaptchaToken string `json:"captcha_token"`
	}{}

	if err := c.ShouldBindJSON(&in); err != nil {
		logger.Warn().Err(err).Msg("invalid request body")
		h.handleError(c, myerrors.NewBadRequestErr("invalid request body"))
		return
	}

	if !utils.IsValidPhoneNumberByCountry(ctx, in.Receiver) {
		logger.Warn().Str("receiver", in.Receiver).Msg("invalid receiver")
		h.handleError(c, myerrors.NewBadRequestErr("invalid phone number"))
		return
	}

	key := fmt.Sprintf("user:%s:send_otp", in.Receiver)
	ok, err := h.limiter.Allow(ctx, key, h.cfg.Service.Security.SendOTPAttempts,
		int(h.cfg.Service.Security.SendOTPBlockTime.Seconds()))
	if err != nil {
		logger.Warn().Err(err).Msg("limiter.Allow error")
		h.handleError(c, myerrors.ErrGeneral)
		return
	}
	if !ok {
		logger.Warn().Msg("sendOTP to phone number is out of limit")
		h.handleError(c, myerrors.NewTooManyRequestsErr(
			"phone number is temporarily blocked due to too many requests"))
		return
	}

	ttl, err := h.service.SendOTP(ctx, in.Receiver)
	if err != nil {
		logger.Error().Err(err).Str("phone", in.Receiver).Msg("service.SendOTP error")
		h.handleError(c, err)
		return
	}

	logger.Debug().Str("receiver", in.Receiver).Msg("OTP sent successfully")
	h.success(c, gin.H{
		"otp_ttl_seconds": ttl,
	})
}

func (h *Handler) confirmOTP(c *gin.Context) {
	ctx := c.Request.Context()
	logger := h.logger.With().
		Ctx(ctx).
		Str("handler", "confirmOTP").
		Logger()

	in := struct {
		Receiver string `json:"receiver" binding:"required"`
		OTP      string `json:"otp" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&in); err != nil {
		logger.Warn().Err(err).Msg("invalid request body")
		h.handleError(c, myerrors.NewBadRequestErr("invalid request body"))
		return
	}

	if !utils.IsValidPhoneNumberByCountry(ctx, in.Receiver) {
		logger.Warn().Str("receiver", in.Receiver).Msg("invalid receiver")
		h.handleError(c, myerrors.NewBadRequestErr("invalid receiver"))
		return
	}

	key := fmt.Sprintf("user:%s:verify_otp", in.Receiver)
	ok, err := h.limiter.Allow(ctx, key, h.cfg.Service.Security.OTPMaxAttempts,
		int(h.cfg.Service.Security.OTPMaxAttemptsBlockTime.Minutes()))
	if err != nil {
		logger.Warn().Err(err).Str("receiver", in.Receiver).Msg("limiter.Allow error")
		h.handleError(c, myerrors.ErrGeneral)
		return
	}
	if !ok {
		logger.Warn().Msg("confirmOTP is out of limit")
		h.handleError(c, myerrors.NewTooManyRequestsErr(
			"receiver is temporarily blocked due to too many requests"))
		return
	}

	response, err := h.service.ConfirmOTP(ctx, in.Receiver, in.OTP)
	if err != nil {
		logger.Error().Err(err).Str("receiver", in.Receiver).Msg("service.ConfirmOTPAndIssueToken error")
		h.handleError(c, err)
		return
	}

	if err = h.limiter.ResetAttempts(ctx, key); err != nil {
		logger.Error().Err(err).Msg("limiter.ResetAttempts error")
	}

	logger.Debug().Str("receiver", in.Receiver).Msg("OTP sent successfully")
	h.success(c, response)
}

func (h *Handler) register(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Структура запроса
	in := struct {
		Email         string `json:"email" binding:"required,email"` // Email обязателен
		Phone         string `json:"phone"`                          // Телефон опционален
		Password      string `json:"password" binding:"required"`
		FullName      string `json:"full_name" binding:"required"`
		// Используем указатель *int, чтобы можно было передать null (для волонтеров)
		InstitutionID *int   `json:"institution_id"`                 
		Role          string `json:"role" binding:"required"`        // 'volunteer' или 'institution'
	}{}

	if err := c.ShouldBindJSON(&in); err != nil {
		h.logger.Warn().Err(err).Msg("invalid register input")
		h.handleError(c, myerrors.NewBadRequestErr("invalid input parameters"))
		return
	}

	// Вызываем сервис
	tokens, err := h.service.Register(ctx, in.Email, in.Phone, in.Password, in.FullName, in.Role, in.InstitutionID)
	if err != nil {
		h.logger.Error().Err(err).Str("email", in.Email).Msg("registration failed")
		h.handleError(c, err)
		return
	}

	h.success(c, tokens)
}

func (h *Handler) login(c *gin.Context) {
	ctx := c.Request.Context()
	logger := h.logger.With().Ctx(ctx).Str("handler", "login").Logger()

	in := struct {
		Email    string `json:"email" binding:"required,email"` // Теперь Email
		Password string `json:"password" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&in); err != nil {
		logger.Warn().Err(err).Msg("invalid login input")
		h.handleError(c, myerrors.NewBadRequestErr("invalid request body"))
		return
	}

	response, err := h.service.Login(ctx, in.Email, in.Password)
	if err != nil {
		logger.Error().Err(err).Str("email", in.Email).Msg("service.Login error")
		h.handleError(c, err)
		return
	}

	logger.Debug().Str("email", in.Email).Msg("login successfully")
	h.success(c, response)
}

// getMe возвращает профиль текущего пользователя
func (h *Handler) getMe(c *gin.Context) {
	// 1. Получаем userID из контекста (его туда положил AuthMiddleware)
	userID, exists := c.Get("userID")
	if !exists {
		h.handleError(c, myerrors.NewUnauthorizedErr("user id not found in context"))
		return
	}

	// 2. Идем в базу
	user, err := h.service.GetUserByID(c.Request.Context(), userID.(int))
	if err != nil {
		h.handleError(c, err)
		return
	}

	// 3. Зачищаем пароль перед отправкой (на всякий случай)
	user.Password = nil

	h.success(c, user)
}
