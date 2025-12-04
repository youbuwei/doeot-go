package repo

import (
	"context"
	"errors"

	"github.com/youbuwei/doeot-go/internal/pay/domain"
	"gorm.io/gorm"
)

// PayModel 是 pay 模块的 GORM 模型。
type PayModel struct {
	ID   int64
	Name string
}

func (PayModel) TableName() string { return "pays" }

// Repo 是基于 GORM 的 domain.Repo 实现。
type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) FindByID(ctx context.Context, id int64) (*domain.Pay, error) {
	var m PayModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrPayNotFound
		}
		return nil, err
	}
	return &domain.Pay{
		ID:   m.ID,
		Name: m.Name,
	}, nil
}

func (r *Repo) Create(ctx context.Context, d *domain.Pay) (*domain.Pay, error) {
	m := PayModel{
		ID:   d.ID,
		Name: d.Name,
	}
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return nil, err
	}
	return &domain.Pay{
		ID:   m.ID,
		Name: m.Name,
	}, nil
}

func (r *Repo) List(ctx context.Context) ([]*domain.Pay, error) {
	var rows []PayModel
	if err := r.db.WithContext(ctx).Find(&rows).Error; err != nil {
		return nil, err
	}
	res := make([]*domain.Pay, 0, len(rows))
	for _, m := range rows {
		res = append(res, &domain.Pay{
			ID:   m.ID,
			Name: m.Name,
		})
	}
	return res, nil
}
