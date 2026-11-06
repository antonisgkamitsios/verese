package app

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/antonisgkamitsios/verese/internal/htmx"
	"github.com/antonisgkamitsios/verese/views/pages"
)

func (app *Application) LogError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	app.Logger.Error(err.Error(), "method", method, "uri", uri)
}

func (app *Application) Render(w http.ResponseWriter, r *http.Request, status int, component templ.Component) {
	w.WriteHeader(status)

	err := component.Render(r.Context(), w)
	if err != nil {
		app.LogError(r, err)
		http.Error(w, http.StatusText(status), status)
	}
}

func (app *Application) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.LogError(r, err)

	w.Header().Add("Vary", "HX-Retarget")

	app.Render(w, r, http.StatusInternalServerError, pages.Error(http.StatusText(http.StatusInternalServerError), "Something went wrong"))
}

func (app *Application) ClientError(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.Header().Add("Vary", "HX-Retarget")

	app.Render(w, r, status, pages.Error(http.StatusText(status), message))
}

func (app *Application) Redirect(w http.ResponseWriter, r *http.Request, location string) {
	if htmx.ContextGetHtmxRequest(r.Context()) {
		w.Header().Add("Hx-Redirect", "/")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, location, http.StatusFound)
}
