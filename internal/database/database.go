package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/antonisgkamitsios/verese/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pressly/goose/v3"
)

type Config struct {
	Dsn         string
	MaxConns    int
	MaxIdleTime time.Duration
}

func OpenDB(config Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(config.Dsn)
	if err != nil {
		return nil, err
	}
	poolConfig.MaxConns = int32(config.MaxConns)
	poolConfig.MaxConnIdleTime = config.MaxIdleTime

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

func Migrate(db *sql.DB) error {
	goose.SetBaseFS(migrations.FS)
	err := goose.SetDialect(string(goose.DialectPostgres))
	if err != nil {
		return err
	}

	return goose.Up(db, ".")
}
