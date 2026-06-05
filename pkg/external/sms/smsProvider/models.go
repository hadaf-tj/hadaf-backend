// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package smsProvider

// SendSMSResponse represents the response from sending SMS
type SendSMSResponse struct {
	Status       string `json:"status"`
	TxnID        string `json:"txn_id"`
	MsgID        string `json:"msg_id"`
	SMSCMsgParts int    `json:"smsc_msg_parts"`
	Timestamp    string `json:"timestamp"`
}

// APIErrorResponse represents an error response from the API
type APIErrorResponse struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Timestamp string `json:"timestamp"`
}

// CheckBalanceResponse represents the response from checking balance
type CheckBalanceResponse struct {
	Balance   string `json:"balance"`
	Timestamp string `json:"timestamp"`
}
