package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handlers) SetupRoutes() http.Handler {
	r := chi.NewRouter()

	r.Method("GET", "/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	r.Get("/", h.HomeHandler)

	return r
}
