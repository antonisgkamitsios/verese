package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/antonisgkamitsios/verese/internal/models"
	"github.com/go-playground/form/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Port int
	Env  string
	Db   struct {
		Dsn         string
		MaxConns    int
		MaxIdleTime time.Duration
	}
}

type Application struct {
	Config      Config
	Logger      *slog.Logger
	Models      models.Models
	FormDecoder *form.Decoder

	// keep track of all background tasks so when graceful shutdown we wait for them to finish
	wg sync.WaitGroup
}

func NewApplication(config Config) (*Application, func()) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(config)
	if err != nil {
		panic(err)
	}
	logger.Info("database connection pool established")

	application := &Application{
		Logger:      logger,
		Config:      config,
		Models:      models.NewModels(db),
		FormDecoder: form.NewDecoder(),
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
	})
}

func openDB(config Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(config.Db.Dsn)
	if err != nil {
		return nil, err
	}
	poolConfig.MaxConns = int32(config.Db.MaxConns)
	poolConfig.MaxConnIdleTime = config.Db.MaxIdleTime

	db, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
