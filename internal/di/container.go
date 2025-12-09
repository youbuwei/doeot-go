package di

import (
	"github.com/youbuwei/doeot-go/internal/app"
	infraDB "github.com/youbuwei/doeot-go/internal/infra/db"
	"github.com/youbuwei/doeot-go/internal/infra/repo"
	httpiface "github.com/youbuwei/doeot-go/internal/interfaces/http"
	httpHandler "github.com/youbuwei/doeot-go/internal/interfaces/http/handler"
	"github.com/youbuwei/doeot-go/pkg/config"
	"go.uber.org/dig"
)

// Build 构建 DI 容器（只包含 HTTP 依赖）
func Build() *dig.Container {
	c := dig.New()

	// 基础设施
	_ = c.Provide(config.Load)
	_ = c.Provide(infraDB.NewGormDB)
	_ = c.Provide(repo.NewOrderRepository)

	// 应用服务
	_ = c.Provide(app.NewOrderService)

	// HTTP handlers
	_ = c.Provide(httpHandler.NewCreateOrderHandler)
	_ = c.Provide(httpHandler.NewGetOrderHandler)
	_ = c.Provide(httpHandler.NewListOrdersHandler)
	_ = c.Provide(httpHandler.NewUpdateOrderHandler)
	_ = c.Provide(httpHandler.NewDeleteOrderHandler)

	// Gin Engine
	_ = c.Provide(httpiface.NewEngine) // *gin.Engine

	return c
}
