package rpcx

import (
    "context"
    "encoding/json"
    "io"
    "net/http"
)

// Request JSON-RPC 2.0 请求
type Request struct {
    JSONRPC string          `json:"jsonrpc"`
    Method  string          `json:"method"`
    Params  json.RawMessage `json:"params"`
    ID      any             `json:"id"`
}

// Response JSON-RPC 2.0 响应
type Response struct {
    JSONRPC string      `json:"jsonrpc"`
    Result  any         `json:"result,omitempty"`
    Error   *Error      `json:"error,omitempty"`
    ID      any         `json:"id"`
}

// Error JSON-RPC 错误结构
type Error struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    any    `json:"data,omitempty"`
}

// Handler JSON-RPC handler 约束
type Handler interface {
    Method() string
    Handle(ctx context.Context, params json.RawMessage) (any, *Error)
}

// Server 简单 JSON-RPC Server 实现
type Server struct {
    handlers map[string]Handler
}

// NewServer 创建 Server 实例
func NewServer() *Server {
    return &Server{
        handlers: make(map[string]Handler),
    }
}

// Register 注册多个 handler
func (s *Server) Register(handlers ...Handler) {
    for _, h := range handlers {
        s.handlers[h.Method()] = h
    }
}

// ServeHTTP 作为 HTTP Handler 暴露
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    var req Request
    if err := json.Unmarshal(body, &req); err != nil {
        writeResponse(w, Response{
            JSONRPC: "2.0",
            Error:   &Error{Code: -32700, Message: "parse error", Data: err.Error()},
            ID:      nil,
        })
        return
    }

    if req.JSONRPC != "2.0" {
        writeResponse(w, Response{
            JSONRPC: "2.0",
            Error:   &Error{Code: -32600, Message: "invalid request", Data: "jsonrpc must be 2.0"},
            ID:      req.ID,
        })
        return
    }

    h, ok := s.handlers[req.Method]
    if !ok {
        writeResponse(w, Response{
            JSONRPC: "2.0",
            Error:   &Error{Code: -32601, Message: "method not found"},
            ID:      req.ID,
        })
        return
    }

    result, rpcErr := h.Handle(r.Context(), req.Params)
    resp := Response{
        JSONRPC: "2.0",
        ID:      req.ID,
    }
    if rpcErr != nil {
        resp.Error = rpcErr
    } else {
        resp.Result = result
    }

    writeResponse(w, resp)
}

func writeResponse(w http.ResponseWriter, resp Response) {
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(resp)
}
