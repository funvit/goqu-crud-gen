package domain

import (
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
