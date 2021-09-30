package goqu_crud_gen

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

// CtxTransaction defines ctx transaction interface which
// allow transfer *sqlx.Tx via context.Context.
type CtxTransaction interface {
	TxFromContext(ctx context.Context) (*sqlx.Tx, error)
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}

// SqlxCtxTran implements CtxTransaction via sqlx.DB.
// Used as default for generated repository.
type SqlxCtxTran struct {
	DB *sqlx.DB
}

func (s *SqlxCtxTran) TxFromContext(ctx context.Context) (*sqlx.Tx, error) {
	return TxFromContext(ctx)
}

func (s *SqlxCtxTran) BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {

	if s.DB == nil {
		return nil, errors.New("repository not initialized, use Connect() first")
	}

	return s.DB.BeginTxx(ctx, opts)
}
