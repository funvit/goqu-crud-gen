package goqu_crud_gen

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
)

type Option func(ds *goqu.SelectDataset)

// WithLockForUpdate option must be used for read methods to lock for update.
func WithLockForUpdate() Option {
	return func(ds *goqu.SelectDataset) {
		orig := *ds
		*ds = *orig.ForUpdate(exp.Wait)
	}
}

type RepositoryOption func(o *RepositoryOpt)

type RepositoryOpt struct {
	TxGetter func(ctx context.Context) (*sqlx.Tx, error)
}

// WithTxGetter option used to set custom transaction getter from context.
func WithTxGetter(fn func(ctx context.Context) (*sqlx.Tx, error)) RepositoryOption {
	return func(o *RepositoryOpt) {
		o.TxGetter = fn
	}
}
