package handler

import (
	"net/http"

	"github.com/antonisgkamitsios/verese/app"
)

type Handlers struct {
	app *app.Application
}

func New(app *app.Application) *Handlers {
	return &Handlers{
		app: app,
	}
}

func (h *Handlers) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	h.app.ClientError(w, r, http.StatusNotFound, "Page not found")
}
