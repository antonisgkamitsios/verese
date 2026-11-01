package handler

import "github.com/antonisgkamitsios/verese/app"

type Handlers struct {
	app *app.Application
}

func New(app *app.Application) *Handlers {
	return &Handlers{
		app: app,
	}
}
