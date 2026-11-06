package main

import (
	"os"
	"time"

	"github.com/antonisgkamitsios/verese/app"
	"github.com/antonisgkamitsios/verese/handler"
	"github.com/antonisgkamitsios/verese/static"
	"github.com/joho/godotenv"
)

func main() {
	var config app.Config

	config.Port = 6969
	config.Env = app.EnvDevelopment

	err := godotenv.Load()
	if err != nil {
		panic("Cannot load .env file")
	}

	config.Db.Dsn = os.Getenv("DB_DSN")
	config.Db.MaxConns = 25
	config.Db.MaxIdleTime = 15 * time.Minute

	config.Smtp.BaseMail = "info@panagia.gr"
	config.Smtp.Host = "sandbox.smtp.mailtrap.io"
	config.Smtp.Username = "8f0d73e77ba826"
	config.Smtp.Password = "848feda30803fb"

	app, close := app.NewApplication(config)
	defer close()
	handlers := handler.New(app)

	static.GenerateHashesForStaticFiles()

	err = app.Serve(handlers.SetupRoutes())

	if err != nil {
		app.Logger.Error(err.Error())
		os.Exit(1)
	}
}
