//go:generate goqu-crud-gen -model Account -table account -dialect mysql -g
package mysql

import (
	"example3/domain"
	"github.com/google/uuid"
)

type Account struct {
	UserId       uuid.UUID `db:"user_id,primary"`
	Login        string    `db:"login"`
	PasswordHash string    `db:"pass"`
}

func (s *Account) Bind(m domain.Account, userId uuid.UUID) error {
	s.UserId = userId
	s.Login = m.Login
	s.PasswordHash = m.PasswordHash

	return nil
}

func (s *Account) ToDomain() (*domain.Account, error) {
	if s == nil {
		return nil, nil
	}

	m := domain.Account{
		Login:        s.Login,
		PasswordHash: s.PasswordHash,
	}

	return &m, nil
}
