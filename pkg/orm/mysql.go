package orm

import (
    "log"
    "time"

    "github.com/youbuwei/doeot-go/pkg/config"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

// NewMySQL creates a *gorm.DB instance with basic pool settings.
func NewMySQL(cfg config.MySQLConfig) *gorm.DB {
    db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to connect mysql: %v", err)
    }

    sqlDB, err := db.DB()
    if err != nil {
        log.Fatalf("failed to get sql.DB: %v", err)
    }

    sqlDB.SetMaxIdleConns(cfg.MaxIdle)
    sqlDB.SetMaxOpenConns(cfg.MaxOpen)
    sqlDB.SetConnMaxLifetime(time.Minute * time.Duration(cfg.MaxLifeMin))

    return db
}
