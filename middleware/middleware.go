package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/antonisgkamitsios/verese/app"
	"github.com/antonisgkamitsios/verese/internal/htmx"
)

type MiddleWare struct {
	app *app.Application
}

func New(app *app.Application) *MiddleWare {
	return &MiddleWare{
		app: app,
	}
}

func (m *MiddleWare) LoggerWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)
		now := time.Now()
		next.ServeHTTP(w, r)
		m.app.Logger.Info("request received",
			"proto", proto,
			"ip", ip,
			"method", method,
			"uri", uri,
			"time", time.Since(now))
	})
}

func (m *MiddleWare) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {

			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				m.app.ServerError(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (m *MiddleWare) DisableCacheInDevMode(next http.Handler) http.Handler {
	if m.app.Config.Env != "development" {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		next.ServeHTTP(w, r)
	})
}

func (m *MiddleWare) CacheStaticFilesInProd(next http.Handler) http.Handler {
	if m.app.Config.Env != "production" {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000")

		next.ServeHTTP(w, r)
	})
}

func (m *MiddleWare) ComesFromHtmx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isHtmxRequest := r.Header.Get("Hx-Request") != ""

		next.ServeHTTP(w, htmx.ContextSetHtmxRequest(r, isHtmxRequest))
	})
}
