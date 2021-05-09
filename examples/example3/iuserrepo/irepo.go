package iuserrepo

import (
	"context"
	"example3/domain"
	"github.com/google/uuid"
)

// UserRepo defines repository interface for domain model domain.User.
type UserRepo interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetForUpdate(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
