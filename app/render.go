package app

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/antonisgkamitsios/verese/internal/htmx"
)

func (app *Application) Render(w http.ResponseWriter, r *http.Request, status int, component templ.Component) {
	w.WriteHeader(status)

	err := component.Render(r.Context(), w)
	if err != nil {
		app.LogError(r, err)
		http.Error(w, http.StatusText(status), status)
	}
}

func (app *Application) Redirect(w http.ResponseWriter, r *http.Request, location string) {
	if htmx.ContextGetHtmxRequest(r.Context()) {
		w.Header().Add("Hx-Redirect", "/")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, location, http.StatusFound)
}
