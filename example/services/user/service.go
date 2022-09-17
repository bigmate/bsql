package user

import (
	"context"
	"errors"
	"github.com/bigmate/bsql/example/models"
	"github.com/bigmate/bsql/example/repository"
	"github.com/bigmate/bsql/example/services/profile"
)

type Service interface {
	Create(ctx context.Context, user models.User) (*models.User, error)
}

type service struct {
	trx     repository.Transactor
	user    repository.User
	profile profile.Service
}

func NewService(trx repository.Transactor, user repository.User, profile profile.Service) Service {
	return &service{
		trx:     trx,
		user:    user,
		profile: profile,
	}
}

func (s *service) Create(ctx context.Context, user models.User) (*models.User, error) {
	ctx, err := s.trx.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer s.trx.Rollback(ctx)

	_, err = s.user.Get(ctx, user.Username)

	if errors.Is(err, repository.ErrNotFound) {
		created, creatErr := s.user.Create(ctx, user)
		if err != nil {
			return nil, creatErr
		}

		_, creatErr = s.profile.Create(ctx, models.Profile{Username: user.Username})
		if creatErr != nil {
			return nil, creatErr
		}

		return created, s.trx.Commit(ctx)
	}

	if err != nil {
		return nil, err
	}

	return nil, errors.New("user exists")
}
