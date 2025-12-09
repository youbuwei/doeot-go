package order

import "context"

// Repository 订单仓储接口（契约）
type Repository interface {
    Create(ctx context.Context, o *Order) error
    FindByID(ctx context.Context, id uint) (*Order, error)
}
