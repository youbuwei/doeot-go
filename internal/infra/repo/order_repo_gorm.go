package repo

import (
    "context"

    "github.com/youbuwei/doeot-go/internal/domain/order"
    "gorm.io/gorm"
)

// OrderRepository GORM 实现的订单仓储
type OrderRepository struct {
    db *gorm.DB
}

// NewOrderRepository 创建订单仓储实现
func NewOrderRepository(db *gorm.DB) order.Repository {
    return &OrderRepository{db: db}
}

// Create 创建订单
func (r *OrderRepository) Create(ctx context.Context, o *order.Order) error {
    return r.db.WithContext(ctx).Create(o).Error
}

// FindByID 根据主键查询订单
func (r *OrderRepository) FindByID(ctx context.Context, id uint) (*order.Order, error) {
    var o order.Order
    if err := r.db.WithContext(ctx).First(&o, id).Error; err != nil {
        return nil, err
    }
    return &o, nil
}
