package goqu_crud_gen

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

type Option func(ds *goqu.SelectDataset)

// WithLockForUpdate option must be used for read methods to lock for update.
func WithLockForUpdate() Option {
	return func(ds *goqu.SelectDataset) {
		orig := *ds
		*ds = *orig.ForUpdate(exp.Wait)
	}
}
