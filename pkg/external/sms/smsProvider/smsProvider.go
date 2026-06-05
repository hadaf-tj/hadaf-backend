// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package smsProvider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"
)

const (
	defaultBaseURL       = "https://api.osonsms.com"
	endpointSendSMS      = "/sendsms_v1.php"
	endpointCheckBalance = "/check_balance.php"
	httpStatusCreated    = 201
	httpStatusOK         = 200
)

type SMSConfig struct {
	APIKey     string
	SenderName string
	Login      string
	BaseURL    string
}

type SMSProvider struct {
	cfg    SMSConfig
	client *httpClient
	logger *zerolog.Logger
}

func NewSMSProvider(cfg SMSConfig, logger *zerolog.Logger) *SMSProvider {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	return &SMSProvider{
		cfg:    cfg,
		client: newHTTPClient(baseURL, cfg.APIKey, cfg.Login),
		logger: logger,
	}
}

func (s *SMSProvider) SendSms(ctx context.Context, phone string, message string, txnID string) error {
	// Format phone number
	formattedPhone, err := FormatPhoneNumber(phone)
	if err != nil {
		s.logger.Error().Ctx(ctx).Err(err).Str("phone", phone).Msg("invalid phone number format")
		return &ValidationError{Message: fmt.Sprintf("invalid phone number: %v", err), Field: "phone"}
	}

	// Build query parameters
	queryParams := map[string]string{
		"from":         s.cfg.SenderName,
		"phone_number": formattedPhone,
		"msg":          message,
		"login":        s.cfg.Login,
		"txn_id":       txnID,
	}

	// Log request (mask phone number for privacy)
	maskedPhone := maskPhoneNumber(formattedPhone)
	s.logger.Info().Ctx(ctx).
		Str("phone", maskedPhone).
		Str("txn_id", txnID).
		Msg("sending SMS")

	// Make API request
	body, statusCode, err := s.client.doRequest(ctx, endpointSendSMS, queryParams)
	if err != nil {
		s.logger.Error().Ctx(ctx).Err(err).
			Str("phone", maskedPhone).
			Str("txn_id", txnID).
			Msg("failed to send SMS")
		return err
	}

	// Parse response
	if statusCode != httpStatusCreated {
		s.logger.Error().Ctx(ctx).
			Int("status_code", statusCode).
			Str("response_body", string(body)).
			Str("txn_id", txnID).
			Msg("SMS API error")
		return parseAPIError(statusCode, body)
	}

	var resp SendSMSResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		s.logger.Error().Ctx(ctx).Err(err).
			Str("txn_id", txnID).
			Msg("failed to parse SMS response")
		return &NetworkError{Message: "failed to parse response", Err: err}
	}

	// Validate response
	if resp.Status != "ok" {
		s.logger.Error().Ctx(ctx).
			Str("status", resp.Status).
			Str("txn_id", txnID).
			Msg("unexpected status in SMS response")
		return &APIError{
			Code:       0,
			Message:    fmt.Sprintf("unexpected status in response: %s", resp.Status),
			HTTPStatus: statusCode,
		}
	}

	s.logger.Info().Ctx(ctx).
		Str("msg_id", resp.MsgID).
		Int("parts", resp.SMSCMsgParts).
		Str("txn_id", txnID).
		Msg("SMS sent successfully")
	return nil
}

// maskPhoneNumber masks phone number for logging (shows only last 4 digits)
func maskPhoneNumber(phone string) string {
	if len(phone) <= 4 {
		return "****"
	}
	return phone[:len(phone)-4] + "****"
}

// BalanceResult represents the result of checking balance
type BalanceResult struct {
	Balance   string
	Timestamp string
}

// CheckBalance checks the account balance
func (s *SMSProvider) CheckBalance(ctx context.Context) (*BalanceResult, error) {
	// Build query parameters
	queryParams := map[string]string{
		"login": s.cfg.Login,
	}

	// Make API request
	body, statusCode, err := s.client.doRequest(ctx, endpointCheckBalance, queryParams)
	if err != nil {
		return nil, err
	}

	// Parse response
	if statusCode != httpStatusOK {
		return nil, parseAPIError(statusCode, body)
	}

	var resp CheckBalanceResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, &NetworkError{Message: "failed to parse response", Err: err}
	}

	return &BalanceResult{
		Balance:   resp.Balance,
		Timestamp: resp.Timestamp,
	}, nil
}
