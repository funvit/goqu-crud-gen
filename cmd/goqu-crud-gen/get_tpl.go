package main

const getTpl = `

// iter iterates other select with specified filter(s).
//
// Can be used in your custom query methods.
func (s *{{ .Repo.Name }}) iter(
	ctx context.Context,
	filter goqu.Expression,
	f func(m {{ .Model.Name }}, stop func()),
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
		var m {{ .Model.Name }}
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

// iterPrimaryKeys iterates other select with specified filter(s).
//
// Can be used in your custom query methods.
func (s *{{ .Repo.Name }}) iterPrimaryKeys(
	ctx context.Context,
	filter goqu.Expression,
	f func(pk interface{}, stop func()),
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
			break
		default:
		}

		err = rows.Scan(&pk)
		if err != nil {
			return fmt.Errorf("row scan error: %w", err)
		}

		f(pk, func() { sigCtxCancel() })
	}

	return nil
}

// each calls wide select.
//
// Can be used in your custom query methods, for example in All.
//
// See also: iter.
func (s *{{ .Repo.Name }}) each(ctx context.Context, f func(m {{ .Model.Name }})) error {

	return s.iter(
		ctx,
		nil,
		func(m {{ .Model.Name }}, _ func()) {
			f(m)
		},
	)
}

// {{"Get"|CRUD}} gets model from database.
//
// Note: returns (nil, nil) if row not found.
func (s *{{ .Repo.Name }}) {{"Get"|CRUD}}(ctx context.Context, id {{.Model.GetPrimaryKeyField.Type}}, opt ...Option) (*{{ .Model.Name }}, error) {

	var r *{{ .Model.Name }}
	err := s.iter(
		ctx,
		s.f.PK().Eq(id),
		func(m {{ .Model.Name }}, stop func()) {
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

func (s *{{ .Repo.Name }}) {{"GetManySlice"|CRUD}}(ctx context.Context, ids []{{.Model.GetPrimaryKeyField.Type}}, opt ...Option) ([]{{ .Model.Name }}, error) {
	items := make([]{{ .Model.Name }}, 0, len(ids))

	err := s.iter(
		ctx,
		s.f.PK().In(ids),
		func(m {{ .Model.Name }}, _ func()) {
			items = append(items, m)
		},
		opt...,
	)
	if err != nil {
		return nil, err
	}

	return items, nil
}
`
