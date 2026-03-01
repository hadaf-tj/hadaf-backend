package sms

import (
	"context"
	"shb/pkg/external/sms/smsProvider"
)

// ISmsAdapter определяет интерфейс для отправки SMS-сообщений.
type ISmsAdapter interface {
	// SendSms отправляет сообщение на указанный номер телефона.
	// txnID - уникальный идентификатор транзакции (обычно ID записи OTP из БД).
	SendSms(ctx context.Context, phone, message, txnID string) error
	// CheckBalance проверяет баланс SMS аккаунта.
	CheckBalance(ctx context.Context) (*smsProvider.BalanceResult, error)
}
