package app

import (
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
