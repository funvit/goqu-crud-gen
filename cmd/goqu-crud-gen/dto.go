package main

const (
	tagOptionAuto    = "auto"
	tagOptionPrimary = "primary"
)

type (
	tplDTO struct {
		GenerateCmd string

		Package string
		Imports []string
		Model   *Model

		Repo Repo

		PrivateCRUD  bool
		WithTranName string
	}
	Repo struct {
		Name    string
		Table   string
		Dialect string
	}
	Model struct {
		MustImport map[string]struct{}
		Name       string
		Fields     []ModelField
	}
	ModelField struct {
		Name    string
		ColName string
		Options []string
		Type    string
	}
)

func (s *ModelField) IsAuto() bool {
	for _, v := range s.Options {
		if v == tagOptionAuto {
			return true
		}
	}

	return false
}

func (s *ModelField) IsPrimary() bool {
	for _, v := range s.Options {
		if v == tagOptionPrimary {
			return true
		}
	}

	return false
}

func (s *Model) GetAutoField() *ModelField {
	for _, f := range s.Fields {
		if f.IsAuto() {
			v := f
			return &v
		}
	}

	return nil
}

func (s *Model) GetPrimaryKeyField() *ModelField {
	for _, f := range s.Fields {
		if f.IsPrimary() {
			v := f
			return &v
		}
	}

	return nil
}

func (s *Model) HasAutoField() bool {
	return s.GetAutoField() != nil
}

func (s *Model) HasPrimaryKeyField() bool {
	return s.GetPrimaryKeyField() != nil
}
