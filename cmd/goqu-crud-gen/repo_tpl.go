package main

const repoTpl = `

type (
	// {{ .Repo.Name }} implements repository for {{ .Model.Name }}.
	{{ .Repo.Name }} struct {
		dsn         string
		db          *sqlx.DB
		dialect     goqu.DialectWrapper
		dialectName string
		options     RepositoryOpt

		// Short for "table".
		t string
		// Short for "table fields", holds repository model fields as goqu exp.IdentifierExpression.
		f {{ .Repo.Name|Private }}Fields
		// Short for "table columns", holds repository model columns as string.
		//
		// Helps to write goqu.UpdateDataset with goqu.Record{}.
		c {{ .Repo.Name|Private }}Columns
	}
	{{ .Repo.Name|Private }}Fields struct {
		{{ range $field := .Model.Fields }}
			{{- $field.Name }} exp.IdentifierExpression
		{{ end }}
	}
	{{ .Repo.Name|Private }}Columns struct {
		{{ range $field := .Model.Fields }}
			{{- $field.Name }} string
		{{ end }}
	}
)

// New{{ .Repo.Name }} returns a new {{ .Repo.Name }}.
//
// Note: do not forget to set max open connections and max lifetime.
func New{{ .Repo.Name }}(dsn string, opt...RepositoryOption) *{{ .Repo.Name }} {
	const t = "{{ .Repo.Table }}"

	s := &{{ .Repo.Name }}{
		dsn:         dsn,
		dialect:     goqu.Dialect("{{ .Repo.Dialect }}"),
		dialectName: "{{ .Repo.Dialect }}",
		t:           t,
		f: {{ .Repo.Name|Private }}Fields{
			{{ range $field := .Model.Fields }}
				{{- $field.Name }}: goqu.C("{{ $field.ColName }}").Table(t),
			{{ end }}
		},
		c: {{ .Repo.Name|Private }}Columns{
			{{ range $field := .Model.Fields }}
				{{- $field.Name }}: "{{ $field.ColName }}",
			{{ end }}
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

// {{ .Repo.Name }}WithInstance returns a new {{ .Repo.Name }} with specified sqlx.DB instance.
func {{ .Repo.Name }}WithInstance(inst *sqlx.DB, opt...RepositoryOption) *{{ .Repo.Name }} {

	const t = "{{ .Repo.Table }}"

	s := &{{ .Repo.Name }}{
		dsn:         "",
		db:          inst,
		dialect:     goqu.Dialect("{{ .Repo.Dialect }}"),
		dialectName: "{{ .Repo.Dialect }}",
		t:           t,
		f: {{ .Repo.Name|Private }}Fields{
			{{ range $field := .Model.Fields }}
				{{- $field.Name }}: goqu.C("{{ $field.ColName }}").Table(t),
			{{ end }}
		},
		c: {{ .Repo.Name|Private }}Columns{
			{{ range $field := .Model.Fields }}
				{{- $field.Name }}: "{{ $field.ColName }}",
			{{ end }}
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
func (s *{{ .Repo.Name|Private }}Fields) PK() exp.IdentifierExpression {
	return s.{{ .Model.GetPrimaryKeyField.Name }}
}

// Connect connects to database instance.
// Must be called after New{{ .Repo.Name }} and before any repo methods.
func (s *{{ .Repo.Name }}) Connect(wait time.Duration) error {

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
func (s *{{ .Repo.Name }}) SetMaxIdleConns(n int) {
	s.db.SetMaxIdleConns(n)
}

// SetMaxOpenConns sets the maximum number of open connections to the database.
//
// See also: sql.SetMaxOpenConns.
func (s *{{ .Repo.Name }}) SetMaxOpenConns(n int) {
	s.db.SetMaxOpenConns(n)
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
//
// See also: sql.SetConnMaxLifetime.
func (s *{{ .Repo.Name }}) SetConnMaxLifetime(d time.Duration) {
	s.db.SetConnMaxLifetime(d)
}

// {{.WithTranName}} wraps function call in transaction.
func (s *{{ .Repo.Name }}) {{.WithTranName}}(ctx context.Context, f func(ctx context.Context) error) error {

	return Transaction(ctx, s.db, s.options.CtxTran, f)
}

// Each query must executed within transaction. This method gets 
// transaction from context, so exists transaction can be used.
func (s *{{ .Repo.Name }}) txFromContext(ctx context.Context) (*sqlx.Tx, error) {

	return s.options.CtxTran.TxFromContext(ctx)
}
`
