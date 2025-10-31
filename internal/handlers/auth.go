package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"shb/pkg/myerrors"
	"shb/pkg/utils"
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
		h.cfg.Service.Security.SendOTPBlockTime)
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

	response, err := h.service.ConfirmOTPAndIssueToken(ctx, in.Receiver, in.OTP)
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

func (h *Handler) register(c *gin.Context) {}

func (h *Handler) login(c *gin.Context) {
	ctx := c.Request.Context()
	logger := h.logger.With().
		Ctx(ctx).
		Str("handler", "login").
		Logger()

	in := struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&in); err != nil {
		logger.Warn().Err(err).Msg("invalid request body")
		h.handleError(c, myerrors.NewBadRequestErr("invalid request body"))
		return
	}

	response, err := h.service.Login(ctx, in.Login, in.Password)
	if err != nil {
		logger.Error().Err(err).Str("phone", in.Login).Msg("service.Login error")
		h.handleError(c, err)
		return
	}

	logger.Debug().Str("login", in.Login).Msg("login successfully")
	h.success(c, response)
}
