package app

import (
	"context"

	"github.com/youbuwei/doeot-go/internal/order/domain"
)

// OrderService 封装了围绕 Order 的业务逻辑。
type OrderService struct {
	repo domain.Repo
}

func NewOrderService(repo domain.Repo) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) Get(ctx context.Context, id int64) (*domain.Order, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *OrderService) Create(ctx context.Context, m *domain.Order) (*domain.Order, error) {
	return s.repo.Create(ctx, m)
}

func (s *OrderService) List(ctx context.Context) ([]*domain.Order, error) {
	return s.repo.List(ctx)
}
