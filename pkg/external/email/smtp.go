package email

import (
	"fmt"
	"net/smtp"
	"os"
)

// ISender интерфейс для отправки (чтобы потом можно было замокать)
type ISender interface {
	SendOTP(to string, code string) error
}

type SmtpSender struct {
	host     string
	port     string
	fromEmail    string
	password string
}

func NewSmtpSender() *SmtpSender {
	// Читаем настройки из ENV. Убедись, что они есть в .env
	return &SmtpSender{
		host:      os.Getenv("SMTP_HOST"),     // например: smtp.gmail.com
		port:      os.Getenv("SMTP_PORT"),     // например: 587
		fromEmail: os.Getenv("SMTP_EMAIL"),    // твоя почта hadaf@gmail.com
		password:  os.Getenv("SMTP_PASSWORD"), // пароль приложения (app password)
	}
}

func (s *SmtpSender) SendOTP(to string, code string) error {
	// Если настройки пустые - просто логируем и выходим (для локальной разработки без инета)
	if s.host == "" || s.password == "" {
		fmt.Printf("[MOCK EMAIL] To: %s, Code: %s\n", to, code)
		return nil
	}

	auth := smtp.PlainAuth("", s.fromEmail, s.password, s.host)
	
	// Простой формат письма
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: Hadaf Verification Code\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/plain; charset=\"utf-8\"\r\n"+
		"\r\n"+
		"Ваш код подтверждения для Hadaf: %s\r\n", to, code))

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	return smtp.SendMail(addr, auth, s.fromEmail, []string{to}, msg)
}