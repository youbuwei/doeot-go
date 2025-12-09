package db

import (
    "github.com/youbuwei/doeot-go/internal/domain/order"
    "github.com/youbuwei/doeot-go/pkg/config"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

// NewGormDB 初始化 GORM（使用 SQLite），并自动迁移 Order 表
func NewGormDB(cfg *config.Config) (*gorm.DB, error) {
    db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    if err := db.AutoMigrate(&order.Order{}); err != nil {
        return nil, err
    }
    return db, nil
}
