// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package smsProvider

import (
	"encoding/json"
	"fmt"
)

// APIError represents an error from the SMS API
type APIError struct {
	Code       int
	Message    string
	HTTPStatus int
	Timestamp  string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("SMS API error [%d]: %s (HTTP %d)", e.Code, e.Message, e.HTTPStatus)
}

// NetworkError represents a network-level error
type NetworkError struct {
	Message string
	Err     error
}

func (e *NetworkError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("network error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("network error: %s", e.Message)
}

func (e *NetworkError) Unwrap() error {
	return e.Err
}

// ValidationError represents an input validation error
type ValidationError struct {
	Message string
	Field   string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error: %s (field: %s)", e.Message, e.Field)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// parseAPIError parses an API error response from JSON body
func parseAPIError(statusCode int, body []byte) error {
	var apiErrResp APIErrorResponse
	if err := json.Unmarshal(body, &apiErrResp); err != nil {
		// If JSON parsing fails, return a generic error
		return &APIError{
			Code:       0,
			Message:    string(body),
			HTTPStatus: statusCode,
		}
	}

	return &APIError{
		Code:       apiErrResp.Code,
		Message:    apiErrResp.Msg,
		HTTPStatus: statusCode,
		Timestamp:  apiErrResp.Timestamp,
	}
}
