package services

import (
	"context"
	"shb/pkg/external/sms/smsProvider"
)

// CheckSMSBalance checks the SMS account balance
func (s *Service) CheckSMSBalance(ctx context.Context) (*smsProvider.BalanceResult, error) {
	return s.sms.CheckBalance(ctx)
}
