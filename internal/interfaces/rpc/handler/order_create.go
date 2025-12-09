package handler

import (
    "context"
    "encoding/json"

    "github.com/youbuwei/doeot-go/internal/app"
    "github.com/youbuwei/doeot-go/pkg/transport/rpcx"
)

// CreateOrderParams JSON-RPC 参数
type CreateOrderParams struct {
    UserID string  `json:"user_id"`
    Amount float64 `json:"amount"`
}

// CreateOrderResult JSON-RPC 结果
type CreateOrderResult struct {
    OrderID uint `json:"order_id"`
}

// CreateOrderRPCHandler JSON-RPC handler
type CreateOrderRPCHandler struct {
    svc app.OrderService
}

// NewCreateOrderRPCHandler 构造函数
func NewCreateOrderRPCHandler(svc app.OrderService) *CreateOrderRPCHandler {
    return &CreateOrderRPCHandler{svc: svc}
}

// Method 返回 JSON-RPC 方法名
func (h *CreateOrderRPCHandler) Method() string {
    return "order.create"
}

// Handle 处理 JSON-RPC 调用
func (h *CreateOrderRPCHandler) Handle(ctx context.Context, raw json.RawMessage) (any, *rpcx.Error) {
    var params CreateOrderParams
    if err := json.Unmarshal(raw, &params); err != nil {
        return nil, &rpcx.Error{Code: -32602, Message: "invalid params", Data: err.Error()}
    }
    if params.UserID == "" || params.Amount <= 0 {
        return nil, &rpcx.Error{Code: -32602, Message: "invalid params", Data: "user_id and amount required"}
    }

    cmd := app.CreateOrderCommand{
        UserID: params.UserID,
        Amount: params.Amount,
    }

    id, err := h.svc.CreateOrder(ctx, cmd)
    if err != nil {
        return nil, &rpcx.Error{Code: 10001, Message: "create order failed", Data: err.Error()}
    }

    return CreateOrderResult{OrderID: id}, nil
}
