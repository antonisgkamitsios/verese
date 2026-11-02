package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handlers) SetupRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(h.RecoverPanic, h.LoggerWrapper)

	r.Group(func(r chi.Router) {
		r.Use(h.CacheStaticFilesInProd, h.DisableCacheInDevMode)
		r.Method("GET", "/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	})
	r.Get("/", h.HomeHandler)

	return r
}
