// Code generated by generator; DO NOT EDIT.

package mysql

import (
	"context"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql" // import is need for proper dialect selection
	"github.com/doug-martin/goqu/v9/exp"
	. "github.com/funvit/goqu-crud-gen"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type (
	// UserPublicFieldsRepo implements repository for UserPublicFields.
	UserPublicFieldsRepo struct {
		dsn         string
		db          *sqlx.DB
		dialect     goqu.DialectWrapper
		dialectName string
		options     RepositoryOpt

		// Short for "table".
		t string
		// Short for "table fields", holds repository model fields as goqu exp.IdentifierExpression.
		f userPublicFieldsRepoFields
		// Short for "table columns", holds repository model columns as string.
		//
		// Helps to write goqu.UpdateDataset with goqu.Record{}.
		c userPublicFieldsRepoColumns
	}
	userPublicFieldsRepoFields struct {
		Id   exp.IdentifierExpression
		Name exp.IdentifierExpression
	}
	userPublicFieldsRepoColumns struct {
		Id   string
		Name string
	}
)

// NewUserPublicFieldsRepo returns a new UserPublicFieldsRepo.
//
// Note: do not forget to set max open connections and max lifetime.
func NewUserPublicFieldsRepo(dsn string, opt ...RepositoryOption) *UserPublicFieldsRepo {
	const t = "user"

	s := &UserPublicFieldsRepo{
		dsn:         dsn,
		dialect:     goqu.Dialect("mysql"),
		dialectName: "mysql",
		t:           t,
		f: userPublicFieldsRepoFields{
			Id:   goqu.C("id").Table(t),
			Name: goqu.C("name").Table(t),
		},
		c: userPublicFieldsRepoColumns{
			Id:   "id",
			Name: "name",
		},
	}
	s.options.CtxTran = &StdCtxTran{
		DB: s.db,
	}

	for _, o := range opt {
		o(&s.options)
	}

	return s
}

// UserPublicFieldsRepoWithInstance returns a new UserPublicFieldsRepo with specified sqlx.DB instance.
func UserPublicFieldsRepoWithInstance(inst *sqlx.DB, opt ...RepositoryOption) *UserPublicFieldsRepo {

	const t = "user"

	s := &UserPublicFieldsRepo{
		dsn:         "",
		db:          inst,
		dialect:     goqu.Dialect("mysql"),
		dialectName: "mysql",
		t:           t,
		f: userPublicFieldsRepoFields{
			Id:   goqu.C("id").Table(t),
			Name: goqu.C("name").Table(t),
		},
		c: userPublicFieldsRepoColumns{
			Id:   "id",
			Name: "name",
		},
	}
	s.options.CtxTran = &StdCtxTran{
		DB: s.db,
	}

	for _, o := range opt {
		o(&s.options)
	}

	return s
}

// PK returns primary key column identifier.
func (s *userPublicFieldsRepoFields) PK() exp.IdentifierExpression {
	return s.Id
}

// Connect connects to database instance.
// Must be called after NewUserPublicFieldsRepo and before any repo methods.
func (s *UserPublicFieldsRepo) Connect(wait time.Duration) error {

	if s.dsn != "" {
		db, err := sqlx.Open(s.dialectName, s.dsn)
		if err != nil {
			return err
		}
		s.db = db
	}

	pCtx, pCancel := context.WithTimeout(context.Background(), wait)
	defer pCancel()

	err := s.db.PingContext(pCtx)
	if err != nil {
		return fmt.Errorf("ping error: %w", err)
	}

	if v, ok := s.options.CtxTran.(*StdCtxTran); ok && v.DB == nil {
		v.DB = s.db
	}

	return nil
}

// SetMaxIdleConns sets the maximum number of connections in the idle
// connection pool.
//
// Note: can helps with stale connections if set (ex: 1 minute).
//
// See also: sql.SetMaxIdleConns.
func (s *UserPublicFieldsRepo) SetMaxIdleConns(n int) {
	s.db.SetMaxIdleConns(n)
}

// SetMaxOpenConns sets the maximum number of open connections to the database.
//
// See also: sql.SetMaxOpenConns.
func (s *UserPublicFieldsRepo) SetMaxOpenConns(n int) {
	s.db.SetMaxOpenConns(n)
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
//
// See also: sql.SetConnMaxLifetime.
func (s *UserPublicFieldsRepo) SetConnMaxLifetime(d time.Duration) {
	s.db.SetConnMaxLifetime(d)
}

// WithTran wraps function call in transaction.
func (s *UserPublicFieldsRepo) WithTran(ctx context.Context, f func(ctx context.Context) error) error {

	return Transaction(ctx, s.db, s.options.CtxTran, f)
}

// Each query must executed within transaction. This method gets
// transaction from context, so exists transaction can be used.
func (s *UserPublicFieldsRepo) txFromContext(ctx context.Context) (*sqlx.Tx, error) {

	return s.options.CtxTran.TxFromContext(ctx)
}

// _Create creates a new row in database by specified model.
//
// If model have "auto" primary key field - it's will be updated in-place.
func (s *UserPublicFieldsRepo) _Create(ctx context.Context, m *UserPublicFields) error {

	tx, err := s.txFromContext(ctx)
	if err != nil {
		return err
	}

	ds := s.dialect.Insert(s.t).Rows(m).Prepared(true)

	q, args, err := ds.ToSQL()
	if err != nil {
		return fmt.Errorf("query builder error: %w", err)
	}

	res, err := tx.Exec(q, args...)
	if err != nil {
		return fmt.Errorf("insert query error: %w", err)
	}
	_ = res

	return nil
}

// iter iterates other select with specified filter(s).
//
// Can be used in your custom query methods.
func (s *UserPublicFieldsRepo) iter(
	ctx context.Context,
	filter goqu.Expression,
	fn func(m UserPublicFields) error,
	opt ...Option,
) error {

	tx, err := s.txFromContext(ctx)
	if err != nil {
		return err
	}

	ds := s.dialect.From(s.t).Prepared(true)

	if filter != nil {
		ds = ds.Where(filter)
	}

	for _, o := range opt {
		o(ds)
	}

	q, args, err := ds.ToSQL()
	if err != nil {
		return fmt.Errorf("query builder error: %w", err)
	}

	sigCtx, sigCtxCancel := context.WithCancel(ctx)
	defer sigCtxCancel()

	rows, err := tx.QueryxContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("select query error: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var m UserPublicFields
		select {
		case <-sigCtx.Done():
			_ = rows.Close()
			return context.Canceled
		default:
		}

		err = rows.StructScan(&m)
		if err != nil {
			return fmt.Errorf("row scan error: %w", err)
		}

		err = fn(m)
		if err != nil {
			sigCtxCancel()
			_ = rows.Close()
			return fmt.Errorf("fn call: %w", err)
		}
	}

	return nil
}

// iterWithOrder iterates other select with specified filter(s) and order.
//
// Can be used in your custom query methods.
func (s *UserPublicFieldsRepo) iterWithOrder(
	ctx context.Context,
	filter goqu.Expression,
	fn func(m UserPublicFields) error,
	order exp.OrderedExpression,
	opt ...Option,
) error {

	tx, err := s.txFromContext(ctx)
	if err != nil {
		return err
	}

	ds := s.dialect.From(s.t).Prepared(true)

	if filter != nil {
		ds = ds.Where(filter)
	}
	if order != nil {
		ds = ds.Order(order)
	}

	for _, o := range opt {
		o(ds)
	}

	q, args, err := ds.ToSQL()
	if err != nil {
		return fmt.Errorf("query builder error: %w", err)
	}

	sigCtx, sigCtxCancel := context.WithCancel(ctx)
	defer sigCtxCancel()

	rows, err := tx.QueryxContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("select query error: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var m UserPublicFields
		select {
		case <-sigCtx.Done():
			_ = rows.Close()
			return context.Canceled
		default:
		}

		err = rows.StructScan(&m)
		if err != nil {
			return fmt.Errorf("row scan error: %w", err)
		}

		err = fn(m)
		if err != nil {
			sigCtxCancel()
			_ = rows.Close()
			return fmt.Errorf("fn call: %w", err)
		}
	}

	return nil
}

// iterPrimaryKeys iterates other select with specified filter(s).
//
// Can be used in your custom query methods.
func (s *UserPublicFieldsRepo) iterPrimaryKeys(
	ctx context.Context,
	filter goqu.Expression,
	fn func(pk interface{}) error,
	opt ...Option,
) error {

	tx, err := s.txFromContext(ctx)
	if err != nil {
		return err
	}

	ds := s.dialect.From(s.t).Prepared(true).Select(s.f.PK())

	if filter != nil {
		ds = ds.Where(filter)
	}

	for _, o := range opt {
		o(ds)
	}

	q, args, err := ds.ToSQL()
	if err != nil {
		return fmt.Errorf("query builder error: %w", err)
	}

	sigCtx, sigCtxCancel := context.WithCancel(ctx)
	defer sigCtxCancel()

	rows, err := tx.QueryxContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("select query error: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var pk interface{}
		select {
		case <-sigCtx.Done():
			_ = rows.Close()
			return context.Canceled
		default:
		}

		err = rows.Scan(&pk)
		if err != nil {
			return fmt.Errorf("row scan error: %w", err)
		}

		err = fn(pk)
		if err != nil {
			sigCtxCancel()
			_ = rows.Close()
			return fmt.Errorf("fn call: %w", err)
		}
	}

	return nil
}

// each calls wide select.
//
// Can be used in your custom query methods, for example in All.
//
// See also: iter.
func (s *UserPublicFieldsRepo) each(ctx context.Context, fn func(m UserPublicFields) error) error {

	return s.iter(
		ctx,
		nil,
		func(m UserPublicFields) error {
			return fn(m)
		},
	)
}

// _Get gets model from database.
//
// Note: returns (nil, nil) if row not found.
func (s *UserPublicFieldsRepo) _Get(ctx context.Context, id uuid.UUID, opt ...Option) (*UserPublicFields, error) {

	var r *UserPublicFields
	err := s.iter(
		ctx,
		s.f.PK().Eq(id),
		func(m UserPublicFields) error {
			// note: expected to be called once.
			r = &m
			return nil
		},
		opt...,
	)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *UserPublicFieldsRepo) _GetManySlice(ctx context.Context, ids []uuid.UUID, opt ...Option) ([]UserPublicFields, error) {
	items := make([]UserPublicFields, 0, len(ids))

	err := s.iter(
		ctx,
		s.f.PK().In(ids),
		func(m UserPublicFields) error {
			items = append(items, m)
			return nil
		},
		opt...,
	)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// _Update updates database row by model.
func (s *UserPublicFieldsRepo) _Update(ctx context.Context, m UserPublicFields) error {

	tx, err := s.txFromContext(ctx)
	if err != nil {
		return err
	}

	ds := s.dialect.Update(s.t).
		Prepared(true).
		Set(m).
		Where(s.f.PK().Eq(m.Id))

	q, args, err := ds.ToSQL()
	if err != nil {
		return fmt.Errorf("query builder error: %w", err)
	}

	_, err = tx.Exec(q, args...)
	if err != nil {
		return fmt.Errorf("update query error: %w", err)
	}

	return nil
}

// _Delete deletes row by id.
//
// Note: returns amount of deleted rows (expected to be max of 1).
//
// See also: _DeleteMany.
func (s *UserPublicFieldsRepo) _Delete(ctx context.Context, id uuid.UUID) (n int64, err error) {

	tx, err := s.txFromContext(ctx)
	if err != nil {
		return 0, err
	}

	ds := s.dialect.Delete(s.t).Where(s.f.PK().Eq(id)).Prepared(true)

	q, args, err := ds.ToSQL()
	if err != nil {
		return 0, fmt.Errorf("query builder error: %w", err)
	}

	res, err := tx.Exec(q, args...)
	if err != nil {
		return 0, fmt.Errorf("delete query error: %w", err)
	}

	return res.RowsAffected()
}

// _DeleteMany deletes rows by ids.
//
// Warning: be careful with large ids arg.
//
// Note: returns amount of deleted rows.
//
// See also: _Delete.
func (s *UserPublicFieldsRepo) _DeleteMany(ctx context.Context, ids []uuid.UUID) (n int64, err error) {

	if len(ids) == 0 {
		// noop
		return 0, nil
	}

	tx, err := s.txFromContext(ctx)
	if err != nil {
		return 0, err
	}

	ds := s.dialect.Delete(s.t).Where(s.f.PK().In(ids)).Prepared(true)

	q, args, err := ds.ToSQL()
	if err != nil {
		return 0, fmt.Errorf("query builder error: %w", err)
	}

	res, err := tx.Exec(q, args...)
	if err != nil {
		return 0, fmt.Errorf("delete query error: %w", err)
	}

	return res.RowsAffected()
}
