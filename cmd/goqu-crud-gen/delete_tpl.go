package main

const deleteTpl = `

// {{"Delete"|CRUD}} deletes row by id.
//
// Note: returns amount of deleted rows (expected to be max of 1).
//
// See also: {{"DeleteMany"|CRUD}}.
func (s *{{.Repo.Name}}) {{"Delete"|CRUD}}(ctx context.Context, id {{.Model.GetPrimaryKeyField.Type}}) (n int64, err error) {

	tx, err := s.txFromContext(ctx)
	if err != nil {
		return 0, err
	}

	ds := s.dialect.Delete(s.t).Where(s.f.PK().Eq(id)).Prepared(true)

	q, args, err := ds.ToSQL()
	if err != nil {
		return 0, fmt.Errorf("query builder: to sql: %w", err)
	}

	res, err := tx.Exec(q, args...)
	if err != nil {
		return 0, fmt.Errorf("tx: exec query: %w", err)
	}

	return res.RowsAffected()
}

// {{"DeleteMany"|CRUD}} deletes rows by ids.
//
// Warning: be careful with large ids arg.
//
// Note: returns amount of deleted rows.
//
// See also: {{"Delete"|CRUD}}.
func (s *{{ .Repo.Name }}) {{"DeleteMany"|CRUD}}(ctx context.Context, ids []{{.Model.GetPrimaryKeyField.Type}}) (n int64, err error) {

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
		return 0, fmt.Errorf("query builder: to sql: %w", err)
	}

	res, err := tx.Exec(q, args...)
	if err != nil {
		return 0, fmt.Errorf("tx: exec: %w", err)
	}

	return res.RowsAffected()
}

`
