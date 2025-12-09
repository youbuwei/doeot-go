package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/youbuwei/doeot-go/internal/app"
	"github.com/youbuwei/doeot-go/pkg/transport/httpx"
)

// ListOrdersResponseItem 列表项 DTO
type ListOrdersResponseItem struct {
	ID     uint    `json:"id"`
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
	Status string  `json:"status"`
}

type ListOrdersHandler struct {
	httpx.Get
	svc app.OrderService
}

func NewListOrdersHandler(svc app.OrderService) *ListOrdersHandler {
	return &ListOrdersHandler{svc: svc}
}

func (h *ListOrdersHandler) Path() string {
	return "/orders"
}

func (h *ListOrdersHandler) Handle(ctx context.Context, c *gin.Context) {
	var userIDPtr *string
	if userID := c.Query("user_id"); userID != "" {
		userIDPtr = &userID
	}

	dtos, err := h.svc.ListOrders(ctx, app.ListOrdersQuery{UserID: userIDPtr})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := make([]ListOrdersResponseItem, 0, len(dtos))
	for _, d := range dtos {
		res = append(res, ListOrdersResponseItem{
			ID:     d.ID,
			UserID: d.UserID,
			Amount: d.Amount,
			Status: d.Status,
		})
	}

	c.JSON(http.StatusOK, res)
}
