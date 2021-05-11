//go:generate goqu-crud-gen -model UserPublicFields -table user -dialect mysql -g -private-crud-methods
package mysql

import (
	"example3/domain"
	"github.com/google/uuid"
)

type UserPublicFields struct {
	Id   uuid.UUID `db:"id,primary"`
	Name string    `db:"name"`
}

func (s *UserPublicFields) Bind(m domain.User) error {

	s.Id = m.Id
	s.Name = m.Name

	return nil
}

func (s *UserPublicFields) ToDomain() (*domain.User, error) {
	if s == nil {
		return nil, nil
	}

	m := domain.User{
		Id:   s.Id,
		Name: s.Name,
	}

	return &m, nil
}
