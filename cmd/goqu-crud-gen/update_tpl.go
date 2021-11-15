package main

const updateTpl = `

// {{"Update"|CRUD}} updates database row by model.
func (s *{{.Repo.Name}}) {{"Update"|CRUD}}(ctx context.Context, m {{.Model.Name}}) error {

	tx, err := s.txFromContext(ctx)
	if err != nil {
		return err
	}

	ds := s.dialect.Update(s.t).
		Prepared(true).
		Set(m).
		Where(s.f.PK().Eq(m.{{.Model.GetPrimaryKeyField.Name}}))

	q, args, err := ds.ToSQL()
	if err != nil {
		return fmt.Errorf("query builder: to sql: %w", err)
	}

	_, err = tx.Exec(q, args...)
	if err != nil {
		return fmt.Errorf("tx: exec: %w", err)
	}

	return nil
}
`
