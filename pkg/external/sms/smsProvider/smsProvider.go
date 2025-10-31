package smsProvider

import (
	"context"
	"shb/pkg/configs"
)

type SMSProvider struct {
	cfg *configs.SMSProvider
}

func NewSMSProvider(cfg *configs.SMSProvider) *SMSProvider {
	return &SMSProvider{cfg: cfg}
}

func (provider *SMSProvider) SendSms(ctx context.Context, phone, message string) error {
	return nil
}
