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

		// short for "table"
		t string
		// short for "table fields"
		f {{ .Repo.Name|Private }}Fields
	}
	{{ .Repo.Name|Private }}Fields struct {
		{{ range $field := .Model.Fields }}
			{{- $field.Name }} exp.IdentifierExpression
		{{ end }}
	}
)

// PK returns primary key column identifier.
func (s *{{ .Repo.Name|Private }}Fields) PK() exp.IdentifierExpression {
	return s.{{ .Model.GetPrimaryKeyField.Name }}
}

// New{{ .Repo.Name }} returns a new {{ .Repo.Name }}.
//
// Note: dont forget to set max open connections and max lifetime.
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
		options: RepositoryOpt{
			TxGetter: GetTxFromContext,
		},
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
	return Transaction(ctx, s.db, f)
}

func (s *{{ .Repo.Name }}) getTxFromContext(ctx context.Context) (*sqlx.Tx, error) {
	return s.options.TxGetter(ctx)
}
`
