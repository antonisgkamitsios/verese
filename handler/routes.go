package handler

import (
	"net/http"

	"github.com/antonisgkamitsios/verese/middleware"
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) SetupRoutes() http.Handler {
	r := chi.NewRouter()

	m := middleware.New(h.app)

	r.Use(m.RecoverPanic, m.LoggerWrapper)

	r.Group(func(r chi.Router) {
		r.Use(m.CacheStaticFilesInProd, m.DisableCacheInDevMode)
		r.Method("GET", "/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	})
	r.Get("/", h.HomeHandler)

	return r
}
