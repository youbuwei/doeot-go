package app

import (
    "context"

    "github.com/youbuwei/doeot-go/internal/domain/order"
)

// CreateOrderCommand 创建订单用例的输入 DTO
type CreateOrderCommand struct {
    UserID string
    Amount float64
}

// OrderService 应用服务接口
type OrderService interface {
    CreateOrder(ctx context.Context, cmd CreateOrderCommand) (uint, error)
}

type orderService struct {
    repo order.Repository
}

// NewOrderService 创建订单应用服务实现
func NewOrderService(repo order.Repository) OrderService {
    return &orderService{repo: repo}
}

// CreateOrder 创建订单
func (s *orderService) CreateOrder(ctx context.Context, cmd CreateOrderCommand) (uint, error) {
    o := &order.Order{
        UserID: cmd.UserID,
        Amount: cmd.Amount,
    }
    if err := s.repo.Create(ctx, o); err != nil {
        return 0, err
    }
    return o.ID, nil
}
