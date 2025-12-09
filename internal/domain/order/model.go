package order

import "time"

// Order 领域实体（同时作为 GORM 模型使用）
type Order struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    string    `gorm:"index;size:64;not null"`
    Amount    float64   `gorm:"not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
