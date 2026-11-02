package handler

import (
	"net/http"
	"time"
)

func (h *Handlers) LoggerWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)
		now := time.Now()
		next.ServeHTTP(w, r)
		h.app.Logger.Info("request received",
			"proto", proto,
			"ip", ip,
			"method", method,
			"uri", uri,
			"time", time.Since(now))
	})
}

func (h *Handlers) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {

			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				// todo: add proper response
				w.Write([]byte("oops"))
				// app.serverError(w, r, fmt.Errorf("%s", err), true)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (h *Handlers) DisableCacheInDevMode(next http.Handler) http.Handler {
	if h.app.Config.Env != "development" {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		next.ServeHTTP(w, r)
	})
}

func (h *Handlers) CacheStaticFilesInProd(next http.Handler) http.Handler {
	if h.app.Config.Env != "production" {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000")

		next.ServeHTTP(w, r)
	})
}
