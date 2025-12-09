package app

import (
	"context"
	"errors"

	"github.com/youbuwei/doeot-go/internal/domain/order"
)

// CreateOrderCommand 创建订单用例输入
type CreateOrderCommand struct {
	UserID string
	Amount float64
}

// UpdateOrderCommand 更新订单用例输入（支持部分更新）
type UpdateOrderCommand struct {
	ID     uint
	Amount *float64
	Status *string
}

// ListOrdersQuery 列表查询参数
type ListOrdersQuery struct {
	UserID *string
}

// OrderDTO 用于对外传递的订单数据
type OrderDTO struct {
	ID     uint
	UserID string
	Amount float64
	Status string
}

// OrderService 应用服务接口（用例集合）
type OrderService interface {
	CreateOrder(ctx context.Context, cmd CreateOrderCommand) (uint, error)
	GetOrder(ctx context.Context, id uint) (*OrderDTO, error)
	ListOrders(ctx context.Context, q ListOrdersQuery) ([]OrderDTO, error)
	UpdateOrder(ctx context.Context, cmd UpdateOrderCommand) error
	DeleteOrder(ctx context.Context, id uint) error
}

type orderService struct {
	repo order.Repository
}

// NewOrderService 创建 OrderService 实现
func NewOrderService(repo order.Repository) OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) CreateOrder(ctx context.Context, cmd CreateOrderCommand) (uint, error) {
	if cmd.UserID == "" {
		return 0, errors.New("user_id required")
	}
	if cmd.Amount <= 0 {
		return 0, errors.New("amount must be > 0")
	}

	o := &order.Order{
		UserID: cmd.UserID,
		Amount: cmd.Amount,
		Status: "created",
	}
	if err := s.repo.Create(ctx, o); err != nil {
		return 0, err
	}
	return o.ID, nil
}

func (s *orderService) GetOrder(ctx context.Context, id uint) (*OrderDTO, error) {
	o, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDTO(o), nil
}

func (s *orderService) ListOrders(ctx context.Context, q ListOrdersQuery) ([]OrderDTO, error) {
	list, err := s.repo.List(ctx, order.Filter{UserID: q.UserID})
	if err != nil {
		return nil, err
	}
	res := make([]OrderDTO, 0, len(list))
	for _, o := range list {
		oCopy := o
		res = append(res, *toDTO(&oCopy))
	}
	return res, nil
}

func (s *orderService) UpdateOrder(ctx context.Context, cmd UpdateOrderCommand) error {
	if cmd.ID == 0 {
		return errors.New("id required")
	}
	o, err := s.repo.Get(ctx, cmd.ID)
	if err != nil {
		return err
	}
	if cmd.Amount != nil {
		if *cmd.Amount <= 0 {
			return errors.New("amount must be > 0")
		}
		o.Amount = *cmd.Amount
	}
	if cmd.Status != nil && *cmd.Status != "" {
		o.Status = *cmd.Status
	}
	return s.repo.Update(ctx, o)
}

func (s *orderService) DeleteOrder(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("id required")
	}
	return s.repo.Delete(ctx, id)
}

func toDTO(o *order.Order) *OrderDTO {
	return &OrderDTO{
		ID:     o.ID,
		UserID: o.UserID,
		Amount: o.Amount,
		Status: o.Status,
	}
}
