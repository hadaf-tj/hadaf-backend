package sms

import "context"

// ISmsAdapter определяет интерфейс для отправки SMS-сообщений.
type ISmsAdapter interface {
	// SendSms отправляет сообщение на указанный номер телефона.
	SendSms(ctx context.Context, phone, message string) error
}
