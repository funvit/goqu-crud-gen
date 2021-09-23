package main

const createTpl = `

// {{ "Create"|CRUD }} creates a new row in database by specified model.
//
// If model have "auto" primary key field - it's will be updated in-place.
func (s *{{.Repo.Name}}) {{"Create"|CRUD}}(ctx context.Context, m *{{.Model.Name}}) error {

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

	{{ if .Model.HasAutoField }}
		m.{{ .Model.GetAutoField.Name }}, err = res.LastInsertId()
		if err != nil {
			return fmt.Errorf("auto id get error: %w", err)
		}
	{{ end }}


	return nil
}
`
