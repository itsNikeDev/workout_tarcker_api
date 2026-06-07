package user

import (
	"context"
	"errors"
	"strings"
)

var ErrNameRequired = errors.New("name is required")

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, name string) (*User, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrNameRequired
	}
	return s.repo.Create(ctx, strings.TrimSpace(name))
}
