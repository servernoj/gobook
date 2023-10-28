package short

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/servernoj/gobook/ch07/sqlx"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type Service struct {
	DB *sqlx.DB
}

func NewService(db *sqlx.DB) *Service {
	return &Service{
		DB: db,
	}
}

func (s *Service) SQLErrorHandling(err error) error {
	if sqliteError, ok := err.(*sqlite.Error); ok {
		switch {
		case sqliteError.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE:
			return fmt.Errorf("UNIQUE constraint violation")
		case sqliteError.Code() == sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
			return fmt.Errorf("PRIMARY KEY constraint violation")
		}
	}
	return nil
}

func (s *Service) Create(ctx context.Context, link Link) error {
	const queryInsert = `insert into links ("short_key","uri") values (?,?)`
	if _, err := s.DB.ExecContext(ctx, queryInsert, link.Key, link.URL); err != nil {
		if err := s.SQLErrorHandling(err); err != nil {
			return err
		}
		return fmt.Errorf("unable to insert value: %w", err)
	}
	return nil
}

func (s *Service) Retrieve(ctx context.Context, key string) (*Link, error) {
	const query = `select l.short_key, l.uri from links as l where l.short_key = ?`
	var link Link
	row := s.DB.QueryRowContext(ctx, query, key)
	if err := row.Scan(&link.Key, &link.URL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &link, nil
}
