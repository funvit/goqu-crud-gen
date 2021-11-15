package main

const createTpl = `

// {{ "Create"|CRUD }} creates a new row in database by specified model.
//
// If model have "auto" primary key field - it's will be updated in-place.
func (s *{{.Repo.Name}}) {{"Create"|CRUD}}(ctx context.Context, m *{{.Model.Name}}) error {

	tx, err := s.txFromContext(ctx)
	if err != nil {
		return err
	}

	ds := s.dialect.Insert(s.t).Rows(m).Prepared(true)

	q, args, err := ds.ToSQL()
	if err != nil {
		return fmt.Errorf("query builder: to sql: %w", err)
	}

	res, err := tx.Exec(q, args...)
	if err != nil {
		return fmt.Errorf("tx: exec: %w", err)
	}
	_ = res

	{{ if .Model.HasAutoField }}
		m.{{ .Model.GetAutoField.Name }}, err = res.LastInsertId()
		if err != nil {
			return fmt.Errorf("tx result: last insert idr: %w", err)
		}
	{{ end }}


	return nil
}
`
