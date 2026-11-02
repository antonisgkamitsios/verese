package main

import (
	"os"

	"github.com/antonisgkamitsios/verese/app"
	"github.com/antonisgkamitsios/verese/handler"
)

func main() {
	var config app.Config

	config.Port = 6969
	config.Env = "development"

	app := app.NewApplication(config)
	handlers := handler.New(app)

	err := app.Serve(handlers.SetupRoutes())

	if err != nil {
		app.Logger.Error(err.Error())
		os.Exit(1)
	}
}
