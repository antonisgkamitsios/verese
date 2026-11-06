package database

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

func SetupTestDb(t *testing.T) *pgxpool.Pool {
	cfg := Config{
		// TODO add test db
		Dsn:         "postgres://verese:pass@localhost/verese?sslmode=disable",
		MaxConns:    25,
		MaxIdleTime: 15 * time.Second,
	}

	db, err := OpenDB(cfg)
	if err != nil {
		t.Fatal(err)
	}

	err = truncateTables(db)
	if err != nil {
		t.Fatal(err)
	}

	err = Migrate(stdlib.OpenDBFromPool(db))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func truncateTables(db *pgxpool.Pool) error {
	rows, err := db.Query(context.Background(),
		`
		SELECT tablename from pg_tables WHERE schemaname = 'public' AND NOT tablename = 'goose_db_version'	
	`)

	if err != nil {
		return err
	}
	defer rows.Close()

	var tables []string

	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			return err
		}
		tables = append(tables, table)
	}

	if rows.Err() != nil {
		return err
	}

	_, err = db.Exec(context.Background(), fmt.Sprintf(`TRUNCATE %s CASCADE`, strings.Join(tables, ", ")))
	if err != nil {
		return err
	}

	return nil
}
