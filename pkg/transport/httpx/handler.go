package httpx

import (
    "context"
    "net/http"

    "github.com/gin-gonic/gin"
)

// Method HTTP 方法类型
type Method string

const (
    MethodGet    Method = http.MethodGet
    MethodPost   Method = http.MethodPost
    MethodPut    Method = http.MethodPut
    MethodDelete Method = http.MethodDelete
)

// Handler 所有 HTTP handler 的统一约束
type Handler interface {
    Method() Method
    Path() string
    Handle(ctx context.Context, c *gin.Context)
}

// Get / Post / Put / Delete 通过嵌入方式约束 Method
type Get struct{}

func (Get) Method() Method { return MethodGet }

type Post struct{}

func (Post) Method() Method { return MethodPost }

type Put struct{}

func (Put) Method() Method { return MethodPut }

type Delete struct{}

func (Delete) Method() Method { return MethodDelete }

// Register 根据 handler 的 Method 和 Path 自动注册到 Gin
func Register(r gin.IRouter, handlers ...Handler) {
    for _, h := range handlers {
        switch h.Method() {
        case MethodGet:
            r.GET(h.Path(), wrap(h))
        case MethodPost:
            r.POST(h.Path(), wrap(h))
        case MethodPut:
            r.PUT(h.Path(), wrap(h))
        case MethodDelete:
            r.DELETE(h.Path(), wrap(h))
        }
    }
}

func wrap(h Handler) gin.HandlerFunc {
    return func(c *gin.Context) {
        h.Handle(c.Request.Context(), c)
    }
}
