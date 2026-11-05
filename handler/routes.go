package handler

import (
	"net/http"

	"github.com/antonisgkamitsios/verese/middleware"
	"github.com/antonisgkamitsios/verese/static"
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) SetupRoutes() http.Handler {
	r := chi.NewRouter()

	m := middleware.New(h.app)
	r.Use(m.RecoverPanic, m.LoggerWrapper, m.ComesFromHtmx, m.Csrf)

	r.NotFound(h.NotFoundHandler)

	r.Group(func(r chi.Router) {
		r.Use(m.CacheStaticFilesInProd, m.DisableCacheInDevMode)
		r.Method("GET", "/static/*", http.StripPrefix("/static", static.StaticFileServer()))
	})
	r.Get("/", h.Home)

	r.Get("/register", h.Register)
	r.Post("/register", h.RegisterPost)

	return r
}
