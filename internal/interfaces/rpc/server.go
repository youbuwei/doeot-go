package rpc

import (
    "net/http"

    "github.com/youbuwei/doeot-go/internal/app"
    "github.com/youbuwei/doeot-go/internal/interfaces/rpc/handler"
    "github.com/youbuwei/doeot-go/pkg/transport/rpcx"
)

// NewHTTPHandler 创建 JSON-RPC HTTP Handler
func NewHTTPHandler(orderSvc app.OrderService) http.Handler {
    server := rpcx.NewServer()
    server.Register(
        handler.NewCreateOrderRPCHandler(orderSvc),
    )

    mux := http.NewServeMux()
    mux.Handle("/rpc", server)

    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        _, _ = w.Write([]byte(`{"status":"ok"}`))
    })

    return mux
}
