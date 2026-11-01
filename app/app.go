package app

import (
	"log"
	"os"
)

type Config struct {
	Port int
	Env  string
}

type Application struct {
	Config Config
	Logger *log.Logger
}

func NewApplication(config Config) *Application {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	return &Application{
		Logger: logger,
		Config: config,
	}
}
