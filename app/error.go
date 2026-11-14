package app

import (
	"net/http"

	"github.com/antonisgkamitsios/verese/views/pages"
)

func (app *Application) LogError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	app.Logger.Error(err.Error(), "method", method, "uri", uri)
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
