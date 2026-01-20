package email

import "context"

// IEmailAdapter определяет интерфейс для отправки email-сообщений.
type IEmailAdapter interface {
	// SendEmail отправляет сообщение на указанный email адрес.
	SendEmail(ctx context.Context, to, subject, body string) error
}
