package application

import (
	"context"
	"time"
)

type Service interface {
	CreateApplication(ctx context.Context, app *Application) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateApplication(ctx context.Context, app *Application) error {
	now := time.Now().UTC()
	if app.CreatedAt.IsZero() {
		app.CreatedAt = now
	}
	app.UpdatedAt = now
	return s.repo.Create(ctx, app)
}
