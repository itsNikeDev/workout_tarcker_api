package exercise

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrNameRequired  = errors.New("name is required")
	ErrAlreadyExists = errors.New("exercise already exists")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, name string, description *string) (*Exercise, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrNameRequired
	}
	return s.repo.Create(ctx, strings.TrimSpace(name), description)
}

func (s *Service) List(ctx context.Context) ([]Exercise, error) {
	return s.repo.List(ctx)
}
