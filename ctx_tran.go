package goqu_crud_gen

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// CtxTransaction defines ctx transaction interface which
// allow transfer *sqlx.Tx via context.Context.
type CtxTransaction interface {
	TxFromContext(ctx context.Context) (*sqlx.Tx, error)
	NewContextWithTx(ctx context.Context, tx *sqlx.Tx) (context.Context, error)
}
