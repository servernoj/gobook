package shortener

import (
	"context"

	"github.com/servernoj/gobook/ch07/short"
	"github.com/servernoj/gobook/ch07/sqlx"
)

type LinkStore interface {
	Create(context.Context, short.Link) error
	Retrieve(context.Context, string) (*short.Link, error)
}

type Service struct {
	LinkStore LinkStore
}

func NewService(db *sqlx.DB) *Service {
	return &Service{
		LinkStore: &short.LinkStore{
			DB: db,
		},
	}
}
