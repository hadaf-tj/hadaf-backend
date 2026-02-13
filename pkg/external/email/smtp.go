package email

import (
	"fmt"
	"net/smtp"
	"os"
)

type ISender interface {
	SendOTP(receiver string, code string) error
}

type SmtpSender struct {
	host     string
	port     string
	email    string
	password string
}

func NewSmtpSender() *SmtpSender {
	return &SmtpSender{
		host:     os.Getenv("SMTP_HOST"),     // smtp.gmail.com
		port:     os.Getenv("SMTP_PORT"),     // 587
		email:    os.Getenv("SMTP_EMAIL"),    // твой_gmail
		password: os.Getenv("SMTP_PASSWORD"), // пароль приложения (App Password)
	}
}

func (s *SmtpSender) SendOTP(receiver string, code string) error {
	if s.email == "" || s.password == "" {
		fmt.Println("[SMTP] Credentials missing in .env, skipping email sending.")
		return nil
	}

	auth := smtp.PlainAuth("", s.email, s.password, s.host)
	addr := s.host + ":" + s.port

	subject := "Subject: Ваш код подтверждения Hadaf\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Здравствуйте!</h2>
			<p>Ваш код подтверждения: <b>%s</b></p>
			<p>Никому не сообщайте этот код.</p>
		</body>
		</html>
	`, code)

	msg := []byte(subject + mime + body)

	if err := smtp.SendMail(addr, auth, s.email, []string{receiver}, msg); err != nil {
		return fmt.Errorf("smtp send error: %w", err)
	}

	return nil
}