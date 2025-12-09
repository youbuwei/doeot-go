package http

import (
	"github.com/gin-gonic/gin"

	"github.com/youbuwei/doeot-go/internal/interfaces/http/handler"
	"github.com/youbuwei/doeot-go/pkg/transport/httpx"
)

// NewEngine 创建 Gin 引擎并注册路由
// 由 DI 容器注入所有 handler 实例。
func NewEngine(
	create *handler.CreateOrderHandler,
	get *handler.GetOrderHandler,
	list *handler.ListOrdersHandler,
	update *handler.UpdateOrderHandler,
	del *handler.DeleteOrderHandler,
) *gin.Engine {
	r := gin.Default()

	httpx.Register(r,
		create,
		get,
		list,
		update,
		del,
	)

	// 简单健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}
