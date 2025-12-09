package http

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"

	"github.com/youbuwei/doeot-go/pkg/transport/httpx"
)

// HandlerIn 集中注入所有实现 httpx.Handler 的 handler
type HandlerIn struct {
	dig.In

	Handlers []httpx.Handler `group:"http.handlers"`
}

// NewEngine 创建 Gin 引擎并注册路由
func NewEngine(in HandlerIn) *gin.Engine {
	r := gin.Default()

	// 把所有 handler 挂到路由上
	if len(in.Handlers) > 0 {
		httpx.Register(r, in.Handlers...)
	}

	// 简单健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}
