package iuserrepo

import (
	"context"

	"example3/domain"
	"github.com/google/uuid"
)

// UserRepo defines repository interface for domain model domain.User.
type UserRepo interface {
	WithTran(ctx context.Context, f func(ctx context.Context) error) error

	Create(ctx context.Context, u domain.User) error

	Get(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetForUpdate(ctx context.Context, id uuid.UUID) (*domain.User, error)

	Delete(ctx context.Context, id uuid.UUID) error
}
