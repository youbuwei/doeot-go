package repo

import (
	"context"
	"errors"

	"github.com/youbuwei/doeot-go/internal/user/domain"
	"gorm.io/gorm"
)

// userModel is the GORM model mapping to the users table.
type userModel struct {
	ID    int64  `gorm:"column:id;primaryKey"`
	Name  string `gorm:"column:name"`
	Age   int    `gorm:"column:age"`
	Role  string `gorm:"column:role"`
	Phone string `gorm:"column:phone"`
}

func (userModel) TableName() string { return "users" }

func (m *userModel) toDomain() *domain.User {
	return &domain.User{
		ID:    m.ID,
		Name:  m.Name,
		Age:   m.Age,
		Role:  m.Role,
		Phone: m.Phone,
	}
}

func fromDomain(u *domain.User) *userModel {
	return &userModel{
		ID:    u.ID,
		Name:  u.Name,
		Age:   u.Age,
		Role:  u.Role,
		Phone: u.Phone,
	}
}

// Repo is a GORM-based implementation of domain.Repo.
type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	var m userModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return m.toDomain(), nil
}

func (r *Repo) List(ctx context.Context) ([]*domain.User, error) {
	var models []userModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}

	res := make([]*domain.User, 0, len(models))
	for i := range models {
		res = append(res, models[i].toDomain())
	}
	return res, nil
}

func (r *Repo) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	m := fromDomain(u)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return m.toDomain(), nil
}
