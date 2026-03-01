package smsProvider

import (
	"fmt"
	"strings"
)

// FormatPhoneNumber converts phone number to 992ХХХХХХХХХ format
func FormatPhoneNumber(phone string) (string, error) {
	// Remove all non-digit characters
	digits := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, phone)

	// Remove leading country code if present
	if strings.HasPrefix(digits, "992") {
		digits = digits[3:]
	}

	// Validate length (should be 9 digits after removing country code)
	if len(digits) != 9 {
		return "", fmt.Errorf("invalid phone number format: expected 9 digits after country code, got %d", len(digits))
	}

	// Return in format 992ХХХХХХХХХ
	return "992" + digits, nil
}
