package htmx

import (
	"context"
	"net/http"
)

type contextKey string

const htmxRequestKey = contextKey("htmxRequest")

func ContextSetHtmxRequest(r *http.Request, isHtmxRequest bool) *http.Request {
	ctx := context.WithValue(r.Context(), htmxRequestKey, isHtmxRequest)

	return r.WithContext(ctx)
}

func ContextGetHtmxRequest(ctx context.Context) bool {
	comesFromHtmx, ok := ctx.Value(htmxRequestKey).(bool)
	if !ok {
		panic("Error: Unable to extract htmxRequest make sure you have called `ContextSetHtmxRequest`, in a middleware perhaps?")
	}

	return comesFromHtmx
}
