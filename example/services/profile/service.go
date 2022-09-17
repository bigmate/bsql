package profile

import (
	"context"
	"errors"
	"github.com/bigmate/bsql/example/models"
	"github.com/bigmate/bsql/example/repository"
)

type Service interface {
	Create(ctx context.Context, profile models.Profile) (*models.Profile, error)
}

type service struct {
	trx     repository.Transactor
	profile repository.Profile
}

func NewService(trx repository.Transactor, profile repository.Profile) Service {
	return &service{
		trx:     trx,
		profile: profile,
	}
}

func (s *service) Create(ctx context.Context, profile models.Profile) (*models.Profile, error) {
	ctx, err := s.trx.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer s.trx.Rollback(ctx)

	_, err = s.profile.Get(ctx, profile.Username)

	if errors.Is(err, repository.ErrNotFound) {
		created, creatErr := s.profile.Create(ctx, profile)
		if creatErr != nil {
			return nil, creatErr
		}

		return created, s.trx.Commit(ctx)
	}

	if err != nil {
		return nil, err
	}

	return nil, errors.New("profile exists")
}
