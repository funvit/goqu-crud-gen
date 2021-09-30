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
	// AccountRepo implements repository for Account.
	AccountRepo struct {
		dsn         string
		db          *sqlx.DB
		dialect     goqu.DialectWrapper
		dialectName string
		options     RepositoryOpt

		// Short for "table".
		t string
		// Short for "table fields", holds repository model fields as goqu exp.IdentifierExpression.
		f accountRepoFields
		// Short for "table columns", holds repository model columns as string.
		//
		// Helps to write goqu.UpdateDataset with goqu.Record{}.
		c accountRepoColumns
	}
	accountRepoFields struct {
		UserId       exp.IdentifierExpression
		Login        exp.IdentifierExpression
		PasswordHash exp.IdentifierExpression
	}
	accountRepoColumns struct {
		UserId       string
		Login        string
		PasswordHash string
	}
)

// PK returns primary key column identifier.
func (s *accountRepoFields) PK() exp.IdentifierExpression {
	return s.UserId
}

// NewAccountRepo returns a new AccountRepo.
//
// Note: dont forget to set max open connections and max lifetime.
func NewAccountRepo(dsn string, opt ...RepositoryOption) *AccountRepo {
	const t = "account"

	s := &AccountRepo{
		dsn:         dsn,
		dialect:     goqu.Dialect("mysql"),
		dialectName: "mysql",
		t:           t,
		f: accountRepoFields{
			UserId:       goqu.C("user_id").Table(t),
			Login:        goqu.C("login").Table(t),
			PasswordHash: goqu.C("pass").Table(t),
		},
		c: accountRepoColumns{
			UserId:       "user_id",
			Login:        "login",
			PasswordHash: "pass",
		},
		options: RepositoryOpt{
			TxGetter: GetTxFromContext,
		},
	}

	for _, o := range opt {
		o(&s.options)
	}

	return s
}

// AccountRepoWithInstance returns a new AccountRepo with specified sqlx.DB instance.
func AccountRepoWithInstance(inst *sqlx.DB, opt ...RepositoryOption) *AccountRepo {

	const t = "account"

	s := &AccountRepo{
		dsn:         "",
		db:          inst,
		dialect:     goqu.Dialect("mysql"),
		dialectName: "mysql",
		t:           t,
		f: accountRepoFields{
			UserId:       goqu.C("user_id").Table(t),
			Login:        goqu.C("login").Table(t),
			PasswordHash: goqu.C("pass").Table(t),
		},
		c: accountRepoColumns{
			UserId:       "user_id",
			Login:        "login",
			PasswordHash: "pass",
		},
		options: RepositoryOpt{
			TxGetter: GetTxFromContext,
		},
	}

	for _, o := range opt {
		o(&s.options)
	}

	return s
}

// Connect connects to database instance.
// Must be called after NewAccountRepo and before any repo methods.
func (s *AccountRepo) Connect(wait time.Duration) error {

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

	return nil
}

// SetMaxIdleConns sets the maximum number of connections in the idle
// connection pool.
//
// Note: can helps with stale connections if set (ex: 1 minute).
//
// See also: sql.SetMaxIdleConns.
func (s *AccountRepo) SetMaxIdleConns(n int) {
	s.db.SetMaxIdleConns(n)
}

// SetMaxOpenConns sets the maximum number of open connections to the database.
//
// See also: sql.SetMaxOpenConns.
func (s *AccountRepo) SetMaxOpenConns(n int) {
	s.db.SetMaxOpenConns(n)
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
//
// See also: sql.SetConnMaxLifetime.
func (s *AccountRepo) SetConnMaxLifetime(d time.Duration) {
	s.db.SetConnMaxLifetime(d)
}

// WithTran wraps function call in transaction.
func (s *AccountRepo) WithTran(ctx context.Context, f func(ctx context.Context) error) error {
	return Transaction(ctx, s.db, f)
}

func (s *AccountRepo) getTxFromContext(ctx context.Context) (*sqlx.Tx, error) {
	return s.options.TxGetter(ctx)
}

// Create creates a new row in database by specified model.
//
// If model have "auto" primary key field - it's will be updated in-place.
func (s *AccountRepo) Create(ctx context.Context, m *Account) error {

	tx, err := s.getTxFromContext(ctx)
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
func (s *AccountRepo) iter(
	ctx context.Context,
	filter goqu.Expression,
	f func(m Account, stop func()),
	opt ...Option,
) error {

	tx, err := s.getTxFromContext(ctx)
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

	// todo: check is it safe to declare var here (not in loop)
	var m Account
	for rows.Next() {
		select {
		case <-sigCtx.Done():
			break
		default:
		}

		err = rows.StructScan(&m)
		if err != nil {
			return fmt.Errorf("row scan error: %w", err)
		}

		f(m, func() { sigCtxCancel() })
	}

	return nil
}

// each calls wide select.
//
// Can be used in your custom query methods, for example in All.
//
// See also: iter.
func (s *AccountRepo) each(ctx context.Context, f func(m Account)) error {

	return s.iter(
		ctx,
		nil,
		func(m Account, _ func()) {
			f(m)
		},
	)
}

// Get gets model from database.
//
// Note: returns (nil, nil) if row not found.
func (s *AccountRepo) Get(ctx context.Context, id uuid.UUID, opt ...Option) (*Account, error) {

	var r *Account
	err := s.iter(
		ctx,
		s.f.PK().Eq(id),
		func(m Account, stop func()) {
			r = &m
			stop()
		},
		opt...,
	)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *AccountRepo) GetManySlice(ctx context.Context, ids []uuid.UUID, opt ...Option) ([]Account, error) {
	items := make([]Account, 0, len(ids))

	err := s.iter(
		ctx,
		s.f.PK().In(ids),
		func(m Account, _ func()) {
			items = append(items, m)
		},
		opt...,
	)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// Update updates database row by model.
func (s *AccountRepo) Update(ctx context.Context, m Account) error {

	tx, err := s.getTxFromContext(ctx)
	if err != nil {
		return err
	}

	ds := s.dialect.Update(s.t).
		Prepared(true).
		Set(m).
		Where(s.f.PK().Eq(m.UserId))

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

// Delete deletes row by id.
//
// Note: returns amount of deleted rows (expected to be max of 1).
//
// See also: DeleteMany.
func (s *AccountRepo) Delete(ctx context.Context, id uuid.UUID) (n int64, err error) {

	tx, err := s.getTxFromContext(ctx)
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

// DeleteMany deletes rows by ids.
//
// Warning: be careful with large ids arg.
//
// Note: returns amount of deleted rows.
//
// See also: Delete.
func (s *AccountRepo) DeleteMany(ctx context.Context, ids []uuid.UUID) (n int64, err error) {

	if len(ids) == 0 {
		// noop
		return 0, nil
	}

	tx, err := s.getTxFromContext(ctx)
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
