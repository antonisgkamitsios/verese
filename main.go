package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/antonisgkamitsios/verese/app"
	"github.com/antonisgkamitsios/verese/handler"
)

func main() {
	var config app.Config

	config.Port = 6969

	app := app.NewApplication(config)
	handlers := handler.New(app)

	server := http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.Port),
		Handler:      handlers.SetupRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server.ListenAndServe()
}
