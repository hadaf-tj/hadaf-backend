package handlers

import (
	"fmt"
	"net/http"
	"os"
	"shb/internal/models"
	"shb/pkg/myerrors"
	"shb/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// setTokenCookies sets httpOnly cookies for access and refresh tokens
func (h *Handler) setTokenCookies(c *gin.Context, tokens *models.TokenResponse) {
	accessMaxAge := int(h.cfg.Security.AccessTokenTTL.Seconds())
	refreshMaxAge := int(h.cfg.Security.RefreshTokenTTL.Seconds())
	isProduction := h.cfg.App.Env == "production"

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		Path:     "/",
		MaxAge:   accessMaxAge,
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/",
		MaxAge:   refreshMaxAge,
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteLaxMode,
	})
}

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

	if !strings.Contains(in.Receiver, "@") && !utils.IsValidPhoneNumberByCountry(ctx, in.Receiver) {
		h.handleError(c, myerrors.NewBadRequestErr("invalid email or phone number"))
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

	in.Receiver = strings.ToLower(strings.TrimSpace(in.Receiver))
	in.OTP = strings.TrimSpace(in.OTP)
	if !strings.Contains(in.Receiver, "@") && !utils.IsValidPhoneNumberByCountry(ctx, in.Receiver) {
		logger.Warn().Str("receiver", in.Receiver).Msg("invalid receiver")
		h.handleError(c, myerrors.NewBadRequestErr("invalid receiver format"))
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

	logger.Debug().Str("receiver", in.Receiver).Msg("OTP confirmed successfully")
	h.setTokenCookies(c, response)
	h.success(c, nil) // Don't send tokens in body
}

func (h *Handler) register(c *gin.Context) {
	ctx := c.Request.Context()

	// Структура запроса
	in := struct {
		Email    string `json:"email" binding:"required,email"` // Email обязателен
		Phone    string `json:"phone"`                          // Телефон опционален
		Password string `json:"password" binding:"required"`
		FullName string `json:"full_name" binding:"required"`
		// Используем указатель *int, чтобы можно было передать null (для волонтеров)
		InstitutionID *int   `json:"institution_id"`
		Role          string `json:"role" binding:"required"` // 'volunteer' или 'institution'
	}{}

	if err := c.ShouldBindJSON(&in); err != nil {
		h.logger.Warn().Err(err).Msg("invalid register input")
		h.handleError(c, myerrors.NewBadRequestErr("invalid input parameters"))
		return
	}

	// C5: Validate password — minimum 8 characters
	if len(in.Password) < 8 {
		h.handleError(c, myerrors.NewBadRequestErr("password must be at least 8 characters"))
		return
	}

	// M4: Validate role — only allow safe values
	if in.Role != "volunteer" && in.Role != "employee" {
		h.handleError(c, myerrors.NewBadRequestErr("invalid role: must be volunteer or employee"))
		return
	}

	in.Email = strings.ToLower(strings.TrimSpace(in.Email))

	// Вызываем сервис
	_, err := h.service.Register(ctx, in.Email, in.Phone, in.Password, in.FullName, in.Role, in.InstitutionID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.success(c, gin.H{
		"message": "verification_required",
		"email":   in.Email,
	})
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

	// Запускаем лимитер ТОЛЬКО если мы не на локалке
	isLocal := os.Getenv("APP_ENV") == "development" || os.Getenv("APP_ENV") == "local"
	if !isLocal {
		// H4: Rate limit login — max 5 attempts per 15 minutes per email
		loginKey := fmt.Sprintf("login:%s", in.Email)
		allowed, err := h.limiter.Allow(ctx, loginKey, 5, 900) // 900 seconds = 15 min
		if err != nil {
			logger.Error().Err(err).Msg("rate limiter error")
		}
		if !allowed {
			h.handleError(c, myerrors.NewTooManyRequestsErr("Слишком много попыток. Повторите позже"))
			return
		}
	}

	response, err := h.service.Login(ctx, in.Email, in.Password)
	if err != nil {
		logger.Error().Err(err).Str("email", in.Email).Msg("service.Login error")
		h.handleError(c, err)
		return
	}

	logger.Debug().Str("email", in.Email).Msg("login successfully")
	h.setTokenCookies(c, response)
	h.success(c, nil) // Don't send tokens in body
}

// refreshTokens handles refresh token rotation
func (h *Handler) refreshTokens(c *gin.Context) {
	ctx := c.Request.Context()
	logger := h.logger.With().Ctx(ctx).Str("handler", "refreshTokens").Logger()

	// Get refresh token from cookie
	cookie, err := c.Request.Cookie("refresh_token")
	if err != nil {
		logger.Warn().Msg("refresh token cookie missing")
		h.handleError(c, myerrors.NewUnauthorizedErr("refresh token missing"))
		return
	}

	response, err := h.service.RefreshTokens(ctx, cookie.Value)
	if err != nil {
		logger.Error().Err(err).Msg("service.RefreshTokens error")
		h.handleError(c, err)
		return
	}

	logger.Debug().Msg("tokens refreshed successfully")
	h.setTokenCookies(c, response)
	h.success(c, nil) // Don't send tokens in body
}

// logout clears httpOnly auth cookies and revokes refresh tokens in DB
func (h *Handler) logout(c *gin.Context) {
	ctx := c.Request.Context()
	userID, exists := c.Get("userID")
	if exists {
		// Revoke all tokens for this user
		if err := h.service.RevokeAllUserRefreshTokens(ctx, userID.(int)); err != nil {
			h.logger.Error().Err(err).Int("userID", userID.(int)).Msg("failed to revoke tokens on logout")
		}
	}

	isProduction := h.cfg.App.Env == "production"

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/api",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteStrictMode,
	})

	h.success(c, gin.H{"message": "logged out"})
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

