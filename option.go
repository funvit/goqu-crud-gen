package goqu_crud_gen

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

type Option func(ds *goqu.SelectDataset)

// WithLockForUpdate option must be used for read methods to lock for update.
func WithLockForUpdate() Option {
	return func(ds *goqu.SelectDataset) {
		*ds = *ds.ForUpdate(exp.Wait)
	}
}

// WithLimit option used to limit select.
func WithLimit(u uint) Option {
	return func(ds *goqu.SelectDataset) {
		*ds = *ds.Limit(u)
	}
}

// WithOrder option used to order select.
func WithOrder(order ...exp.OrderedExpression) Option {
	return func(ds *goqu.SelectDataset) {
		*ds = *ds.Order(order...)
	}
}

// WithFilter option used to filter select.
func WithFilter(exp ...exp.Expression) Option {
	return func(ds *goqu.SelectDataset) {
		*ds = *ds.Where(exp...)
	}
}

type RepositoryOption func(o *RepositoryOpt)

type RepositoryOpt struct {
	CtxTran CtxTransaction
}

// WithCtxTran option used to set custom transaction getter from context.
func WithCtxTran(ct CtxTransaction) RepositoryOption {
	return func(o *RepositoryOpt) {
		o.CtxTran = ct
	}
}
