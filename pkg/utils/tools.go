// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package utils

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"regexp"
	"shb/pkg/constants"
	"strings"
)

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func GenerateOTP(length int) (string, error) {
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length {
		return "", err
	}
	for i := range b {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b), nil
}

func IsValidPhoneNumberByCountry(ctx context.Context, phone string) bool {
	countryCode, ok := ctx.Value(constants.CountryCodeKey).(string)
	if !ok || countryCode == "" {
		countryCode = "TJ"
	}

	switch countryCode {
	case "TJ":
		return isValidTajikPhone(phone)
	case "RU":
		return isValidRussianPhone(phone)
	// Additional country codes can be added here.
	default:
		return false
	}
}

func isValidTajikPhone(phone string) bool {
	// Expected format: 992 + 9 digits = 12 characters total.
	if len(phone) != 12 || !strings.HasPrefix(phone, "992") {
		return false
	}
	return isDigits(phone)
}

func isValidRussianPhone(phone string) bool {
	// Expected format: 7 + 10 digits = 11 characters total.
	if strings.HasPrefix(phone, "8") && len(phone) == 11 {
		phone = "7" + phone[1:]
	}
	if len(phone) != 11 || !strings.HasPrefix(phone, "7") {
		return false
	}
	return isDigits(phone)
}

var digitsRegexp = regexp.MustCompile(`^\d+$`)

func isDigits(s string) bool {
	return digitsRegexp.MatchString(s)
}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
