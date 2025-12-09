package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/youbuwei/doeot-go/internal/app"
	"github.com/youbuwei/doeot-go/pkg/transport/httpx"
)

// GetOrderResponse GET /orders/:id 响应 DTO
type GetOrderResponse struct {
	ID     uint    `json:"id"`
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
	Status string  `json:"status"`
}

type GetOrderHandler struct {
	httpx.Get
	svc app.OrderService
}

func NewGetOrderHandler(svc app.OrderService) *GetOrderHandler {
	return &GetOrderHandler{svc: svc}
}

func (h *GetOrderHandler) Path() string {
	return "/orders/:id"
}

func (h *GetOrderHandler) Handle(ctx context.Context, c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	dto, err := h.svc.GetOrder(ctx, uint(id64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, GetOrderResponse{
		ID:     dto.ID,
		UserID: dto.UserID,
		Amount: dto.Amount,
		Status: dto.Status,
	})
}
