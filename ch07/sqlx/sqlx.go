package sqlx

import (
	"context"
	"database/sql"
	"fmt"

	_ "embed"

	_ "modernc.org/sqlite"
)

const dsn = "ch07.db"

//go:embed schema.sql
var schema string

type DB struct {
	*sql.DB
}

func Dial() (*DB, error) {
	ctx := context.Background()
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening DB driver: %w", err)
	}
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("pinging DB: %w", err)
	}
	if _, err := db.ExecContext(ctx, schema); err != nil {
		return nil, fmt.Errorf("applying schema: %w", err)
	}
	return &DB{DB: db}, nil
}
