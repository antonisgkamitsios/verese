package reqctx

import (
	"context"
	"net/http"

	"github.com/antonisgkamitsios/verese/internal/models"
)

type contextKey string

const (
	htmxKey contextKey = "htmx"
	csrfKey contextKey = "csrf"
	userKey contextKey = "user"
)

func RequestWithHtmx(r *http.Request, isHtmxRequest bool) *http.Request {
	ctx := context.WithValue(r.Context(), htmxKey, isHtmxRequest)

	return r.WithContext(ctx)
}

func IsHtmx(ctx context.Context) bool {
	comesFromHtmx, ok := ctx.Value(htmxKey).(bool)
	if !ok {
		panic("Error: Unable to extract htmxRequest make sure you have called `WithHtmx`, in a middleware perhaps?")
	}

	return comesFromHtmx
}

func RequestWithCsrfToken(r *http.Request, csrfToken string) *http.Request {
	ctx := context.WithValue(r.Context(), csrfKey, csrfToken)

	return r.WithContext(ctx)
}

func CSRF(ctx context.Context) string {
	if token, ok := ctx.Value(csrfKey).(string); ok {
		return token
	}

	return ""
}

func RequestWithUser(r *http.Request, user *models.User) *http.Request {
	ctx := context.WithValue(r.Context(), userKey, user)

	return r.WithContext(ctx)
}

func User(ctx context.Context) *models.User {

	user, ok := ctx.Value(userKey).(*models.User)
	if !ok {
		panic("missing user value, make sure you have wrapped called `RequestWithUser`")
	}

	return user
}
