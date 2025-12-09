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

func (r *OrderRepository) Create(ctx context.Context, o *order.Order) error {
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *OrderRepository) Get(ctx context.Context, id uint) (*order.Order, error) {
	var o order.Order
	if err := r.db.WithContext(ctx).First(&o, id).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *OrderRepository) List(ctx context.Context, f order.Filter) ([]order.Order, error) {
	db := r.db.WithContext(ctx)
	if f.UserID != nil {
		db = db.Where("user_id = ?", *f.UserID)
	}
	var list []order.Order
	if err := db.Order("id DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *OrderRepository) Update(ctx context.Context, o *order.Order) error {
	return r.db.WithContext(ctx).Save(o).Error
}

func (r *OrderRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&order.Order{}, id).Error
}
