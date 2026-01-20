package smtpEmail

import (
	"context"
	"fmt"
	"net/smtp"
	"shb/internal/configs"
	"strconv"
)

type SMTPEmail struct {
	cfg *configs.SMTPConfig
}

func NewSMTPEmail(cfg *configs.SMTPConfig) *SMTPEmail {
	return &SMTPEmail{cfg: cfg}
}

func (s *SMTPEmail) SendEmail(ctx context.Context, to, subject, body string) error {
	if s.cfg.Host == "" || s.cfg.Port == "" {
		return fmt.Errorf("SMTP host and port must be configured")
	}

	port, err := strconv.Atoi(s.cfg.Port)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %w", err)
	}

	addr := fmt.Sprintf("%s:%d", s.cfg.Host, port)
	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)

	from := s.cfg.FromEmail
	if s.cfg.FromName != "" {
		from = fmt.Sprintf("%s <%s>", s.cfg.FromName, s.cfg.FromEmail)
	}

	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n", from, to, subject, body))

	err = smtp.SendMail(addr, auth, s.cfg.FromEmail, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
