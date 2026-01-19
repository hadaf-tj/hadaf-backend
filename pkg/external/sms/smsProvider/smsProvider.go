package smsProvider

import (
	"context"
	"fmt"
)

type SMSConfig struct {
	APIKey     string
	SenderName string
}

type SMSProvider struct {
	cfg SMSConfig
}

func NewSMSProvider(cfg SMSConfig) *SMSProvider {
	return &SMSProvider{cfg: cfg}
}

func (s *SMSProvider) SendSms(ctx context.Context, phone string, message string) error {
	// Mock implementation
	fmt.Printf("[SMS MOCK] To: %s, Body: %s, Sender: %s\n", phone, message, s.cfg.SenderName)
	return nil
}