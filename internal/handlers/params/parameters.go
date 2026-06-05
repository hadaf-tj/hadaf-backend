// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package params

import (
	"context"
	"shb/pkg/constants"
)

func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(constants.RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}
