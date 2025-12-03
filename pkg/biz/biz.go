package biz

import (
    "context"
    "encoding/json"
)

// Context is the abstraction exposed to business handlers.
// It hides the underlying HTTP / RPC frameworks.
type Context interface {
    RequestContext() context.Context
    RequestID() string

    Bind(out any) error
    JSON(status int, body any) error
    Result(data any, err error) error
}

// HandlerFunc is the function signature used by HTTP handlers.
type HandlerFunc func(Context) error

// Router is an abstract HTTP router used by modules to register routes.
type Router interface {
    GET(path string, h HandlerFunc, opts ...RouteOption)
    POST(path string, h HandlerFunc, opts ...RouteOption)
    PUT(path string, h HandlerFunc, opts ...RouteOption)
    DELETE(path string, h HandlerFunc, opts ...RouteOption)
}

// RPCHandlerFunc is the function signature used by RPC handlers.
type RPCHandlerFunc func(Context, json.RawMessage) (any, error)

// RPCRouter is an abstract RPC router used by modules to register JSON-RPC methods.
type RPCRouter interface {
    Handle(method string, h RPCHandlerFunc, opts ...RouteOption)
}

// RouteMeta stores metadata driven by annotations (auth/tags/etc).
type RouteMeta struct {
    Auth   string
    Tags   []string
    BizTag string
}

// RouteOption mutates RouteMeta.
type RouteOption func(*RouteMeta)

func WithAuth(auth string) RouteOption {
    return func(m *RouteMeta) {
        m.Auth = auth
    }
}

func WithTags(tags ...string) RouteOption {
    return func(m *RouteMeta) {
        m.Tags = append(m.Tags, tags...)
    }
}

func WithBizTag(tag string) RouteOption {
    return func(m *RouteMeta) {
        m.BizTag = tag
    }
}

// Module is implemented by each business module (user/order/...).
type Module interface {
    Name() string
    RegisterHTTP(r Router)
    RegisterRPC(r RPCRouter)
}
