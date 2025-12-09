package di

import (
    "github.com/youbuwei/doeot-go/internal/app"
    infraDB "github.com/youbuwei/doeot-go/internal/infra/db"
    "github.com/youbuwei/doeot-go/internal/infra/repo"
    "github.com/youbuwei/doeot-go/pkg/config"
    "go.uber.org/dig"
)

// Build 构建 DI 容器
func Build() *dig.Container {
    c := dig.New()

    _ = c.Provide(config.Load)
    _ = c.Provide(infraDB.NewGormDB)
    _ = c.Provide(repo.NewOrderRepository)
    _ = c.Provide(app.NewOrderService)

    return c
}
