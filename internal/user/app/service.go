package app

import (
	"context"

	"github.com/youbuwei/doeot-go/internal/user/domain"
)

// UserService holds business logic around User.
type UserService struct {
	repo domain.Repo
}

func NewUserService(repo domain.Repo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) GetUserList(ctx context.Context) ([]*domain.User, error) {
	return s.repo.List(ctx)
}

func (s *UserService) CreateUser(ctx context.Context, u *domain.User) (*domain.User, error) {
	return s.repo.Create(ctx, u)
}
