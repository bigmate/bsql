package repository

import (
	"context"
	"github.com/bigmate/bsql"
	"github.com/bigmate/bsql/example/models"
)

type User interface {
	Get(ctx context.Context, username string) (*models.User, error)
	Create(ctx context.Context, user models.User) (*models.User, error)
}

type user struct {
	db bsql.Factory
}

func NewUserRepository(db bsql.Factory) User {
	return &user{db: db}
}

func (u *user) Get(ctx context.Context, username string) (*models.User, error) {
	fetched := &models.User{}
	query := "SELECT * FROM user WHERE username = $"

	if err := u.db.FromContext(ctx).GetContext(ctx, fetched, query, username); err != nil {
		return nil, err
	}

	return fetched, nil
}

func (u *user) Create(ctx context.Context, user models.User) (*models.User, error) {
	created := &models.User{}
	query := "INSERT INTO user(username) VALUES ($) RETURNING *"

	if err := u.db.FromContext(ctx).GetContext(ctx, created, query, user.CreatedAt); err != nil {
		return nil, err
	}

	return created, nil
}

type Profile interface {
	Get(ctx context.Context, username string) (*models.Profile, error)
	Create(ctx context.Context, profile models.Profile) (*models.Profile, error)
}

type profile struct {
	db bsql.Factory
}

func (p *profile) Get(ctx context.Context, username string) (*models.Profile, error) {
	fetched := &models.Profile{}
	query := "SELECT * FROM profile WHERE username = $"

	if err := p.db.FromContext(ctx).GetContext(ctx, fetched, query, username); err != nil {
		return nil, err
	}

	return fetched, nil
}

func (p *profile) Create(ctx context.Context, profile models.Profile) (*models.Profile, error) {
	created := &models.Profile{}
	query := "INSERT INTO user VALUES ($, $, $, $) RETURNING *"
	err := p.db.FromContext(ctx).GetContext(
		ctx,
		created,
		query,
		profile.Username,
		profile.FirstName,
		profile.LastName,
		profile.Country,
	)

	if err != nil {
		return nil, err
	}

	return created, nil
}

func NewProfileRepository(db bsql.Factory) Profile {
	return &profile{db: db}
}
