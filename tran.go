package goqu_crud_gen

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ctxTxKeyType string

const ctxTxKey = ctxTxKeyType("repo.sqlx.tx")

var ErrNoTranInContext = fmt.Errorf("no transaction in context")

// TxFromContext gets started transaction from context.
func TxFromContext(ctx context.Context) (*sqlx.Tx, error) {
	v := ctx.Value(ctxTxKey)
	if v == nil {
		return nil, ErrNoTranInContext
	}

	tx, ok := v.(*sqlx.Tx)
	if !ok {
		return nil, fmt.Errorf("cant cast tx from context to *sqlx.Tx")
	}

	return tx, nil
}

// Transaction calls func in transaction. Rollbacks if function return error.
//
// Example:
//
//    err := Transaction(ctx, ct, func(ctx context.Context) error {
//        m, err := ...
//        if err != nil {
//            return err
//        }
//
//        // other db operation
//        return ...
//    })
//
// Special method for generated repositories.
func Transaction(ctx context.Context, ct CtxTransaction, f func(ctx context.Context) error) error {
	// if tx already in ctx - use it
	tx, err := ct.TxFromContext(ctx)
	if err == nil && tx != nil {
		return f(ctx)
	}

	// new tx
	tx, err = ct.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	txCtx := context.WithValue(ctx, ctxTxKey, tx)

	defer func() {
		if err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil {
				err = fmt.Errorf("tran rollback error: %w", rbErr)
			}
			return
		}
		err = tx.Commit()
		if err != nil {
			err = fmt.Errorf("tran commit error: %w", err)
			return
		}
	}()

	return f(txCtx)
}
