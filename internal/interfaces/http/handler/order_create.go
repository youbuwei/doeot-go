package handler

import (
    "context"
    "net/http"

    "github.com/gin-gonic/gin"

    "github.com/youbuwei/doeot-go/internal/app"
    "github.com/youbuwei/doeot-go/pkg/transport/httpx"
)

// CreateOrderRequest 创建订单请求 DTO
type CreateOrderRequest struct {
    UserID string  `json:"user_id" binding:"required"`
    Amount float64 `json:"amount"  binding:"required,gt=0"`
}

// CreateOrderResponse 创建订单响应 DTO
type CreateOrderResponse struct {
    OrderID uint `json:"order_id"`
}

// CreateOrderHandler 单接口单文件的 HTTP handler
type CreateOrderHandler struct {
    httpx.Post           // 继承 POST 约束
    svc app.OrderService // 依赖应用服务
}

// NewCreateOrderHandler 构造函数
func NewCreateOrderHandler(svc app.OrderService) *CreateOrderHandler {
    return &CreateOrderHandler{svc: svc}
}

// Path 路径定义
func (h *CreateOrderHandler) Path() string {
    return "/orders"
}

// Handle 处理请求
func (h *CreateOrderHandler) Handle(_ context.Context, c *gin.Context) {
    var req CreateOrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    cmd := app.CreateOrderCommand{
        UserID: req.UserID,
        Amount: req.Amount,
    }

    id, err := h.svc.CreateOrder(c.Request.Context(), cmd)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, CreateOrderResponse{OrderID: id})
}
