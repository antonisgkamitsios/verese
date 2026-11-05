package csrf

import "context"

type contextKey string

const CSRFContextKey = contextKey("csrf")

func GetCSRF(ctx context.Context) string {
	if token, ok := ctx.Value(CSRFContextKey).(string); ok {
		return token
	}

	return ""
}
