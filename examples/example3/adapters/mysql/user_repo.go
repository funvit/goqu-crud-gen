package mysql

import (
	"context"
	"example3/domain"
	"fmt"
	. "github.com/funvit/goqu-crud-gen"
	"github.com/google/uuid"
	"time"
)

type UserRepo struct {
	*UserPublicFieldsRepo
	accRepo *AccountRepo
}

func BuildUser(p UserPublicFields, a *Account) (*domain.User, error) {
	dp, err := p.ToDomain()
	if err != nil {
		return nil, err
	}

	da, err := a.ToDomain()
	if err != nil {
		return nil, err
	}

	m := domain.User{
		Id:      dp.Id,
		Name:    dp.Name,
		Account: da,
	}

	return &m, nil
}

// NewUserRepo returns a new UserRepo.
func NewUserRepo(dsn string) *UserRepo {
	return &UserRepo{
		UserPublicFieldsRepo: NewUserPublicFieldsRepo(dsn),
		accRepo:              NewAccountRepo(dsn),
	}
}

func (s UserRepo) Connect(wait time.Duration) error {
	err := s.UserPublicFieldsRepo.Connect(wait)
	if err != nil {
		return err
	}

	// use same conn poll for account repo.
	s.accRepo.db = s.UserPublicFieldsRepo.db

	return nil
}

func (s *UserRepo) get(ctx context.Context, id uuid.UUID, opt ...Option) (*domain.User, error) {

	// fixme: second Get flaky: bad conn

	p, err := s._Get(ctx, id, opt...)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}

	a, err := s.accRepo.Get(ctx, id, opt...)
	if err != nil {
		return nil, err
	}

	return BuildUser(*p, a)
}

func (s *UserRepo) Create(ctx context.Context, u domain.User) error {

	var upf UserPublicFields
	err := upf.Bind(u)
	if err != nil {
		return fmt.Errorf("model convert error: %w", err)
	}

	err = s._Create(ctx, &upf)
	if err != nil {
		return err
	}

	var acc Account
	err = acc.Bind(*u.Account, u.Id)
	if err != nil {
		return fmt.Errorf("model convert error: %w", err)
	}

	err = s.accRepo.Create(ctx, &acc)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserRepo) Get(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.get(ctx, id)
}

func (s *UserRepo) GetForUpdate(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.get(ctx, id, WithLockForUpdate())
}

func (s *UserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s._Delete(ctx, id)
	if err != nil {
		return nil
	}

	_, err = s.accRepo.Delete(ctx, id)
	if err != nil {
		return nil
	}

	return nil
}
