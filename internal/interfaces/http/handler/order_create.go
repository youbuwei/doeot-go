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

// CreateOrderHandler POST /orders
type CreateOrderHandler struct {
	httpx.Post
	svc app.OrderService
}

func NewCreateOrderHandler(svc app.OrderService) *CreateOrderHandler {
	return &CreateOrderHandler{svc: svc}
}

func (h *CreateOrderHandler) Path() string {
	return "/orders"
}

func (h *CreateOrderHandler) Handle(_ context.Context, c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.svc.CreateOrder(c.Request.Context(), app.CreateOrderCommand{
		UserID: req.UserID,
		Amount: req.Amount,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, CreateOrderResponse{OrderID: id})
}
