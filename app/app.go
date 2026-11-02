package app

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
)

type Config struct {
	Port int
	Env  string
}

type Application struct {
	Config Config
	Logger *slog.Logger

	// keep track of all background tasks so when graceful shutdown we wait for them to finish
	wg sync.WaitGroup
}

func NewApplication(config Config) *Application {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return &Application{
		Logger: logger,
		Config: config,
	}
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
