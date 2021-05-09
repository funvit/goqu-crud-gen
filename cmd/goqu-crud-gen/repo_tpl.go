package main

const repoTpl = `

type (
	// {{ .Repo.Name }} implements repository for {{ .Model.Name }}.
	{{ .Repo.Name }} struct {
		dsn         string
		db          *sqlx.DB
		dialect     goqu.DialectWrapper
		dialectName string

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
func New{{ .Repo.Name }}(dsn string) *{{ .Repo.Name }} {
	const t = "{{ .Repo.Table }}"

	return &{{ .Repo.Name }}{
		dsn:         dsn,
		dialect:     goqu.Dialect("{{ .Repo.Dialect }}"),
		dialectName: "{{ .Repo.Dialect }}",
		t:           t,
		f: {{ .Repo.Name|Private }}Fields{
			{{ range $field := .Model.Fields }}
				{{- $field.Name }}: goqu.C("{{ $field.ColName }}").Table(t),
			{{ end }}
		},
	}
}

// Connect connects to database instance.
// Must be called after New{{ .Repo.Name }} and before any repo methods.
func (s *{{ .Repo.Name }}) Connect(wait time.Duration) error {
	db, err := sqlx.Open(s.dialectName, s.dsn)
	if err != nil {
		return err
	}

	pCtx, pCancel := context.WithTimeout(context.Background(), wait)
	defer pCancel()
	err = db.PingContext(pCtx)
	if err != nil {
		return fmt.Errorf("ping error: %w", err)
	}

	s.db = db

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

// WithTran wraps function call in transaction.
func (s *{{ .Repo.Name }}) WithTran(ctx context.Context, f func(ctx context.Context) error) error {
	return Transaction(ctx, s.db, f)
}
`
