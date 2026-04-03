package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_SendOTP_sms(t *testing.T) {
	ctx := context.Background()
	svc, d := newTestService(t)
	d.Repo.On("SaveOTP", ctx, mock.AnythingOfType("*models.OTP")).Return(1, nil)
	d.SMS.On("SendSms", mock.Anything, "+1234567890", mock.AnythingOfType("string"), "1").Return(nil).Maybe()
	ttl, err := svc.SendOTP(ctx, "+1234567890")
	require.NoError(t, err)
	require.Equal(t, 300, ttl)
}

func TestService_SendOTP_email(t *testing.T) {
	ctx := context.Background()
	svc, d := newTestService(t)
	d.Repo.On("SaveOTP", ctx, mock.AnythingOfType("*models.OTP")).Return(2, nil)
	d.Email.On("SendEmail", mock.Anything, "u@example.com", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Maybe()
	_, err := svc.SendOTP(ctx, "u@example.com")
	require.NoError(t, err)
}
