package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/youbuwei/doeot-go/internal/app"
	"github.com/youbuwei/doeot-go/pkg/transport/httpx"
)

type DeleteOrderHandler struct {
	httpx.Delete
	svc app.OrderService
}

func NewDeleteOrderHandler(svc app.OrderService) *DeleteOrderHandler {
	return &DeleteOrderHandler{svc: svc}
}

func (h *DeleteOrderHandler) Path() string {
	return "/orders/:id"
}

func (h *DeleteOrderHandler) Handle(ctx context.Context, c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.svc.DeleteOrder(ctx, uint(id64)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
