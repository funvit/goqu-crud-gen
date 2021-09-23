package domain

import (
	"fmt"

	"github.com/google/uuid"
)

type (
	// User defines domain model.
	User struct {
		Id      uuid.UUID
		Name    string
		Account *Account
	}
	Account struct {
		Login        string
		PasswordHash string
	}
)

func (u *User) String() string {
	return fmt.Sprintf(
		"User(Id:%q Name:%q Account:%s)",
		u.Id.String(),
		u.Name,
		u.Account.String(),
	)
}

func (a *Account) String() string {
	if a == nil {
		return "nil"
	}
	return fmt.Sprintf("Account(Login:%q)", a.Login)
}
