package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/youbuwei/doeot-go/internal/app"
	"github.com/youbuwei/doeot-go/pkg/transport/httpx"
)

// UpdateOrderRequest 更新订单请求 DTO（支持部分更新）
type UpdateOrderRequest struct {
	Amount *float64 `json:"amount,omitempty"`
	Status *string  `json:"status,omitempty"`
}

type UpdateOrderHandler struct {
	httpx.Put
	svc app.OrderService
}

func NewUpdateOrderHandler(svc app.OrderService) *UpdateOrderHandler {
	return &UpdateOrderHandler{svc: svc}
}

func (h *UpdateOrderHandler) Path() string {
	return "/orders/:id"
}

func (h *UpdateOrderHandler) Handle(ctx context.Context, c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := app.UpdateOrderCommand{
		ID:     uint(id64),
		Amount: req.Amount,
		Status: req.Status,
	}

	if err := h.svc.UpdateOrder(ctx, cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
