package order

import (
	"context"
	"time"
)

// Order 领域实体，同时作为 GORM 模型
type Order struct {
	ID        uint    `gorm:"primaryKey"`
	UserID    string  `gorm:"size:64;not null;index"`
	Amount    float64 `gorm:"not null"`
	Status    string  `gorm:"size:32;not null;default:'created'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Filter 查询过滤条件
type Filter struct {
	UserID *string
}

// Repository 仓储契约（属于领域层）
type Repository interface {
	Create(ctx context.Context, o *Order) error
	Get(ctx context.Context, id uint) (*Order, error)
	List(ctx context.Context, f Filter) ([]Order, error)
	Update(ctx context.Context, o *Order) error
	Delete(ctx context.Context, id uint) error
}
