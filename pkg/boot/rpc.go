package boot

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"

	"github.com/youbuwei/doeot-go/pkg/biz"
	"github.com/youbuwei/doeot-go/pkg/errs"
)

// JSON-RPC basic request/response types.

type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      json.RawMessage `json:"id"`
}

type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  any             `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
	ID      json.RawMessage `json:"id"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type rpcServer struct {
	addr     string
	handlers map[string]biz.RPCHandlerFunc
}

func newRPCServer(addr string) *rpcServer {
	return &rpcServer{
		addr:     addr,
		handlers: make(map[string]biz.RPCHandlerFunc),
	}
}

// rpcRouter adapts rpcServer to biz.RPCRouter.
type rpcRouter struct {
	srv *rpcServer
}

func (r *rpcRouter) Handle(method string, h biz.RPCHandlerFunc, opts ...biz.RouteOption) {
	// meta could be used for logging/metrics/ACL; ignored in demo.
	_ = opts
	r.srv.handlers[method] = h
}

func (a *App) runRPC() error {
	srv := newRPCServer(a.cfg.RPC.Addr)
	router := &rpcRouter{srv: srv}

	for _, m := range a.modules {
		m.RegisterRPC(router)
	}

	return srv.start()
}

func (s *rpcServer) start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handle)

	httpSrv := &http.Server{
		Addr:    s.addr,
		Handler: mux,
	}

	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	return httpSrv.Serve(ln)
}

func (s *rpcServer) handle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req rpcRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeRPCError(w, req.ID, -32700, "parse error")
		return
	}

	h, ok := s.handlers[req.Method]
	if !ok {
		writeRPCError(w, req.ID, -32601, "method not found")
		return
	}

	// Build a minimal biz.Context for RPC.
	ctx := &rpcContext{ctx: r.Context()}

	result, err := h(ctx, req.Params)
	if err != nil {
		var e *errs.Error
		if errors.As(err, &e) {
			writeRPCError(w, req.ID, mapErrCode(e.Code), e.Msg)
			return
		}
		writeRPCError(w, req.ID, -32000, "internal error")
		return
	}

	resp := rpcResponse{
		JSONRPC: "2.0",
		Result:  result,
		ID:      req.ID,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func writeRPCError(w http.ResponseWriter, id json.RawMessage, code int, msg string) {
	resp := rpcResponse{
		JSONRPC: "2.0",
		Error: &rpcError{
			Code:    code,
			Message: msg,
		},
		ID: id,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func mapErrCode(c errs.Code) int {
	switch c {
	case errs.CodeBadRequest:
		return -32602
	case errs.CodeNotFound:
		return -32004
	case errs.CodeInternal:
		fallthrough
	default:
		return -32000
	}
}

// rpcContext is a minimal implementation of biz.Context for RPC calls.
// Bind/JSON/Result are not used in this demo (wrappers work directly with params/result).
type rpcContext struct {
	ctx context.Context
}

func (c *rpcContext) RequestContext() context.Context {
	return c.ctx
}

func (c *rpcContext) RequestID() string {
	return ""
}

func (c *rpcContext) Bind(out any) error {
	return errors.New("Bind not supported for RPC context; use params decoding in wrapper")
}

func (c *rpcContext) JSON(status int, body any) error {
	_ = status
	_ = body
	return errors.New("JSON not supported for RPC context")
}

func (c *rpcContext) Result(data any, err error) error {
	_ = data
	_ = err
	return errors.New("Result not supported for RPC context")
}
