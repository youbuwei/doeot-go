package app

import (
	"context"

	"github.com/youbuwei/doeot-go/internal/pay/domain"
)

// PayService 封装了围绕 Pay 的业务逻辑。
type PayService struct {
	repo domain.Repo
}

func NewPayService(repo domain.Repo) *PayService {
	return &PayService{repo: repo}
}

func (s *PayService) Get(ctx context.Context, id int64) (*domain.Pay, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *PayService) Create(ctx context.Context, m *domain.Pay) (*domain.Pay, error) {
	return s.repo.Create(ctx, m)
}

func (s *PayService) List(ctx context.Context) ([]*domain.Pay, error) {
	return s.repo.List(ctx)
}
