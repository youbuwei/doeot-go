package repo

import (
	"context"
	"errors"

	"github.com/youbuwei/doeot-go/internal/order/domain"
	"gorm.io/gorm"
)

// OrderModel 是 order 模块的 GORM 模型。
type OrderModel struct {
	ID   int64
	Name string
}

func (OrderModel) TableName() string { return "orders" }

// Repo 是基于 GORM 的 domain.Repo 实现。
type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) FindByID(ctx context.Context, id int64) (*domain.Order, error) {
	var m OrderModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrOrderNotFound
		}
		return nil, err
	}
	return &domain.Order{
		ID:   m.ID,
		Name: m.Name,
	}, nil
}

func (r *Repo) Create(ctx context.Context, d *domain.Order) (*domain.Order, error) {
	m := OrderModel{
		ID:   d.ID,
		Name: d.Name,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return nil, err
	}
	return &domain.Order{
		ID:   m.ID,
		Name: m.Name,
	}, nil
}

func (r *Repo) List(ctx context.Context) ([]*domain.Order, error) {
	var rows []OrderModel
	if err := r.db.WithContext(ctx).Find(&rows).Error; err != nil {
		return nil, err
	}
	res := make([]*domain.Order, 0, len(rows))
	for _, m := range rows {
		res = append(res, &domain.Order{
			ID:   m.ID,
			Name: m.Name,
		})
	}
	return res, nil
}
