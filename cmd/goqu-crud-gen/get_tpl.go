package main

const getTpl = `

// iter iterates other select.
//
// Can be used in your custom query methods.
//
// Filters, limit or order can be set via opts.
func (s *{{ .Repo.Name }}) iter(
	ctx context.Context,
	fn func(m {{ .Model.Name }}) error,
	opt ...Option,
) error {

	tx, err := s.txFromContext(ctx)
	if err != nil {
		return err
	}

	ds := s.dialect.From(s.t).Prepared(true)

	for _, o := range opt {
		o(ds)
	}

	q, args, err := ds.ToSQL()
	if err != nil {
		return fmt.Errorf("query builder: to sql: %w", err)
	}

	sigCtx, sigCtxCancel := context.WithCancel(ctx)
	defer sigCtxCancel()

	rows, err := tx.QueryxContext(sigCtx, q, args...)
	if err != nil {
		return fmt.Errorf("tx: query rows: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var m {{ .Model.Name }}

		err = ctx.Err()
		if err != nil {
			return fmt.Errorf("rows: next: %w", err)
		}

		err = rows.StructScan(&m)
		if err != nil {
			return fmt.Errorf("rows struct scan: %w", err)
		}

		err = fn(m)
		if err != nil {
			sigCtxCancel()
			return fmt.Errorf("fn call: %w", err)
		}
	}

	return nil
}

// iterPrimaryKeys iterates other select with specified filter(s).
//
// Can be used in your custom query methods.
//
// Filters, limit or order can be set via opts.
func (s *{{ .Repo.Name }}) iterPrimaryKeys(
	ctx context.Context,
	fn func(pk interface{}) error,
	opt ...Option,
) error {

	tx, err := s.txFromContext(ctx)
	if err != nil {
		return err
	}

	ds := s.dialect.From(s.t).Prepared(true).Select(s.f.PK())

	for _, o := range opt {
		o(ds)
	}

	q, args, err := ds.ToSQL()
	if err != nil {
		return fmt.Errorf("query builder: to sql: %w", err)
	}

	sigCtx, sigCtxCancel := context.WithCancel(ctx)
	defer sigCtxCancel()

	rows, err := tx.QueryxContext(sigCtx, q, args...)
	if err != nil {
		return fmt.Errorf("tx: query rows: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var pk interface{}

		err = ctx.Err()
		if err != nil {
			return fmt.Errorf("rows: next: %w", err)
		}

		err = rows.Scan(&pk)
		if err != nil {
			return fmt.Errorf("row scan: %w", err)
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

// {{"Get"|CRUD}} gets model from database.
//
// Note: returns (nil, nil) if row not found.
//
// See also: {{"GetForUpdate"|CRUD}}.
func (s *{{ .Repo.Name }}) {{"Get"|CRUD}}(ctx context.Context, id {{.Model.GetPrimaryKeyField.Type}}) (*{{ .Model.Name }}, error) {

	var r *{{ .Model.Name }}

	opts := []Option{
		WithFilter(s.f.PK().Eq(id)),
	}

	err := s.iter(
		ctx,
		func(m {{ .Model.Name }}) error {
			// note: expected to be called once.
			r = &m
			return nil
		},
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// {{"GetForUpdate"|CRUD}} gets model from database for update (i.e. locks row).
//
// Note: returns (nil, nil) if row not found.
//
// See also: {{"Get"|CRUD}}.
func (s *{{ .Repo.Name }}) {{"GetForUpdate"|CRUD}}(ctx context.Context, id {{.Model.GetPrimaryKeyField.Type}}) (*{{ .Model.Name }}, error) {

	var r *{{ .Model.Name }}

	opts := []Option{
		WithFilter(s.f.PK().Eq(id)),
		WithLockForUpdate(),
	}

	err := s.iter(
		ctx,
		func(m {{ .Model.Name }}) error {
			// note: expected to be called once.
			r = &m
			return nil
		},
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// {{"GetMany"|CRUD}} gets models from database.
//
// See also: {{"GetManyForUpdate"|CRUD}}.
func (s *{{ .Repo.Name }}) {{"GetMany"|CRUD}}(ctx context.Context, ids []{{.Model.GetPrimaryKeyField.Type}}) ([]{{ .Model.Name }}, error) {

	items := make([]{{ .Model.Name }}, 0, len(ids))

	opts := []Option{
		WithFilter(s.f.PK().In(ids)),
	}

	err := s.iter(
		ctx,
		func(m {{ .Model.Name }}) error {
			items = append(items, m)
			return nil
		},
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// {{"GetManyForUpdate"|CRUD}} gets models from database for update (i.e. locks rows).
//
// See also: {{"GetMany"|CRUD}}.
func (s *{{ .Repo.Name }}) {{"GetManyForUpdate"|CRUD}}(ctx context.Context, ids []{{.Model.GetPrimaryKeyField.Type}}) ([]{{ .Model.Name }}, error) {

	items := make([]{{ .Model.Name }}, 0, len(ids))

	opts := []Option{
		WithFilter(s.f.PK().In(ids)),
		WithLockForUpdate(),
	}

	err := s.iter(
		ctx,
		func(m {{ .Model.Name }}) error {
			items = append(items, m)
			return nil
		},
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return items, nil
}
`
