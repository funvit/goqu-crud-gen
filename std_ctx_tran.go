package goqu_crud_gen

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

type ctxTxKeyType string

const ctxTxKey = ctxTxKeyType("repo.sqlx.tx")

var ErrNoTranInContext = errors.New("no transaction in context")

// StdCtxTran implements CtxTransaction via sqlx.DB.
// Used as default for generated repository.
type StdCtxTran struct {
	DB *sqlx.DB
}

func (s *StdCtxTran) TxFromContext(ctx context.Context) (*sqlx.Tx, error) {
	return TxFromContext(ctx)
}

func (s *StdCtxTran) NewContextWithTx(ctx context.Context, tx *sqlx.Tx) (context.Context, error) {
	return context.WithValue(ctx, ctxTxKey, tx), nil
}
