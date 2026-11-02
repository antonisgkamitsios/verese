package handler

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/antonisgkamitsios/verese/app"
	"github.com/antonisgkamitsios/verese/views/pages"
)

type Handlers struct {
	app *app.Application
}

func New(app *app.Application) *Handlers {
	return &Handlers{
		app: app,
	}
}

func (h *Handlers) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	h.app.Logger.Error(err.Error(), "method", method, "uri", uri)
}

func (h *Handlers) serverError(w http.ResponseWriter, r *http.Request, err error) {
	h.logError(r, err)

	w.Header().Add("Vary", "HX-Retarget")

	h.render(w, r, http.StatusInternalServerError, pages.Error(http.StatusText(http.StatusInternalServerError), "Something went wrong"))
}

func (h *Handlers) clientError(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.Header().Add("Vary", "HX-Retarget")

	h.render(w, r, status, pages.Error(http.StatusText(status), message))
}

func (h *Handlers) render(w http.ResponseWriter, r *http.Request, status int, component templ.Component) {
	w.WriteHeader(status)

	err := component.Render(r.Context(), w)
	if err != nil {
		h.logError(r, err)
		http.Error(w, http.StatusText(status), status)
	}
}

func (h *Handlers) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	h.clientError(w, r, http.StatusNotFound, "Page not found")
}
