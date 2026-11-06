package app

import (
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/antonisgkamitsios/verese/internal/database"
	"github.com/antonisgkamitsios/verese/internal/mailer"
	"github.com/antonisgkamitsios/verese/internal/models"
	"github.com/go-playground/form/v4"
)

const (
	EnvDevelopment = "developlmen"
	EnvProduction  = "production"
	EnvTesting     = "testing"
)

type Config struct {
	Port int
	Env  string
	Db   database.Config
	Smtp mailer.Config
}

type Application struct {
	Config      Config
	Logger      *slog.Logger
	Models      models.Models
	FormDecoder *form.Decoder
	Mailer      mailer.Mailer

	// keep track of all background tasks so when graceful shutdown we wait for them to finish
	wg sync.WaitGroup
}

func NewApplication(config Config) (*Application, func()) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := database.OpenDB(config.Db)
	if err != nil {
		panic(err)
	}
	logger.Info("database connection pool established")

	application := &Application{
		Logger:      logger,
		Config:      config,
		Models:      models.NewModels(db),
		FormDecoder: form.NewDecoder(),
		Mailer:      mailer.New(config.Smtp),
	}

	return application, db.Close
}

func (app *Application) Background(fn func()) {
	app.wg.Go(func() {

		defer func() {
			if err := recover(); err != nil {
				app.Logger.Error(fmt.Sprintf("%v", err))
			}
		}()

		fn()
	})
}
