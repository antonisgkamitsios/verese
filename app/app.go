package app

import (
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/alexedwards/scs/v2"
	"github.com/antonisgkamitsios/verese/internal/database"
	"github.com/antonisgkamitsios/verese/internal/mailer"
	"github.com/antonisgkamitsios/verese/internal/models"
	"github.com/antonisgkamitsios/verese/internal/session"
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
	FormDecoder *form.Decoder
	Logger      *slog.Logger
	Mailer      mailer.Mailer
	Models      models.Models
	Session     *scs.SessionManager

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
		Config:      config,
		FormDecoder: form.NewDecoder(),
		Logger:      logger,
		Mailer:      mailer.New(config.Smtp),
		Models:      models.NewModels(db),
		Session:     session.NewSessionManager(db),
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
