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
