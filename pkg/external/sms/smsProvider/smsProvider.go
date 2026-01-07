package smsProvider

import (
	"context"
	"shb/internal/configs"
)

type SMSProvider struct {
	cfg *configs.SMSConfig
}

func NewSMSProvider(cfg *configs.SMSConfig) *SMSProvider {
	return &SMSProvider{cfg: cfg}
}

func (provider *SMSProvider) SendSms(ctx context.Context, phone, message string) error {
	return nil
}
