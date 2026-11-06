package app

import (
	"io"
	"log/slog"
	"testing"

	"github.com/antonisgkamitsios/verese/internal/database"
	"github.com/antonisgkamitsios/verese/internal/models"
	"github.com/go-playground/form/v4"
)

func NewTestApplication(t *testing.T) *Application {
	testCfg := Config{
		Port: 69420,
		Env:  EnvTesting,
	}
	db := database.SetupTestDb(t)

	application := &Application{
		Logger:      slog.New(slog.NewTextHandler(io.Discard, nil)),
		Config:      testCfg,
		Models:      models.NewModels(db),
		FormDecoder: form.NewDecoder(),
	}

	return application
}
