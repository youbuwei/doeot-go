package http

import (
    "github.com/gin-gonic/gin"

    "github.com/youbuwei/doeot-go/internal/app"
    "github.com/youbuwei/doeot-go/internal/interfaces/http/handler"
    "github.com/youbuwei/doeot-go/pkg/transport/httpx"
)

// NewEngine 创建 Gin 引擎并注册路由
func NewEngine(orderSvc app.OrderService) *gin.Engine {
    r := gin.Default()

    // 注册业务 handler
    httpx.Register(r,
        handler.NewCreateOrderHandler(orderSvc),
    )

    // 一个简单的健康检查
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    return r
}
