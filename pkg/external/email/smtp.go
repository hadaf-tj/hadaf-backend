package email

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
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

	// 1. Очищаем email от случайных пробелов и скрытых символов
	cleanReceiver := strings.TrimSpace(receiver)

	// 2. Строго по стандарту: используем \r\n и добавляем обязательные To и From
	headers := "From: " + s.email + "\r\n" +
		"To: " + cleanReceiver + "\r\n" +
		"Subject: Ваш код подтверждения Hadaf\r\n" +
		"MIME-version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" // <--- Эта пустая строка критически важна! Она отделяет заголовки от тела

	body := fmt.Sprintf(`
        <html>
        <body>
            <h2>Здравствуйте!</h2>
            <p>Ваш код подтверждения: <b>%s</b></p>
            <p>Никому не сообщайте этот код.</p>
        </body>
        </html>
    `, code)

	// Склеиваем правильные заголовки и тело
	msg := []byte(headers + body)

	// 3. Отправляем на очищенный адрес
	if err := smtp.SendMail(addr, auth, s.email, []string{cleanReceiver}, msg); err != nil {
		return fmt.Errorf("smtp send error: %w", err)
	}

	return nil
}