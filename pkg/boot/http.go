package boot

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/youbuwei/doeot-go/pkg/biz"
	"github.com/youbuwei/doeot-go/pkg/errs"
)

func (a *App) runHTTP() error {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	router := &echoRouter{e: e}

	for _, m := range a.modules {
		m.RegisterHTTP(router)
	}

	return e.Start(a.cfg.HTTP.Addr)
}

// echoRouter adapts echo.Echo to biz.Router.
type echoRouter struct {
	e *echo.Echo
}

func (r *echoRouter) wrap(h biz.HandlerFunc, meta *biz.RouteMeta) echo.HandlerFunc {
	// meta can be used to attach middleware like auth/metrics; omitted in demo.
	_ = meta

	return func(c echo.Context) error {
		ctx := newEchoContext(c)
		return h(ctx)
	}
}

func buildRouteMeta(opts []biz.RouteOption) *biz.RouteMeta {
	m := &biz.RouteMeta{}
	for _, o := range opts {
		o(m)
	}
	return m
}

func (r *echoRouter) GET(path string, h biz.HandlerFunc, opts ...biz.RouteOption) {
	meta := buildRouteMeta(opts)
	r.e.GET(path, r.wrap(h, meta))
}

func (r *echoRouter) POST(path string, h biz.HandlerFunc, opts ...biz.RouteOption) {
	meta := buildRouteMeta(opts)
	r.e.POST(path, r.wrap(h, meta))
}

func (r *echoRouter) PUT(path string, h biz.HandlerFunc, opts ...biz.RouteOption) {
	meta := buildRouteMeta(opts)
	r.e.PUT(path, r.wrap(h, meta))
}

func (r *echoRouter) DELETE(path string, h biz.HandlerFunc, opts ...biz.RouteOption) {
	meta := buildRouteMeta(opts)
	r.e.DELETE(path, r.wrap(h, meta))
}

// echoContext implements biz.Context on top of echo.Context.
type echoContext struct {
	c echo.Context
}

func newEchoContext(c echo.Context) *echoContext {
	return &echoContext{c: c}
}

func (ctx *echoContext) RequestContext() context.Context {
	return ctx.c.Request().Context()
}

func (ctx *echoContext) RequestID() string {
	id := ctx.c.Request().Header.Get("X-Request-ID")
	if id == "" {
		id = ctx.c.Response().Header().Get("X-Request-ID")
	}
	return id
}

// Bind supports a tiny subset of binding rules:
//   - JSON body (via echo.Bind)
//   - `path:"name"` tags from URL parameters
func (ctx *echoContext) Bind(out any) error {
	// First let echo try JSON/query/form binding.
	if err := ctx.c.Bind(out); err != nil {
		// ignore here; we still try manual binding below
		_ = err
	}

	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("Bind expects non-nil pointer")
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("Bind expects pointer to struct")
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fv := v.Field(i)
		if !fv.CanSet() {
			continue
		}

		if pathName := field.Tag.Get("path"); pathName != "" {
			val := ctx.c.Param(pathName)
			if val == "" {
				continue
			}
			if err := setValue(fv, val); err != nil {
				return err
			}
		}
	}

	return nil
}

func setValue(fv reflect.Value, s string) error {
	switch fv.Kind() {
	case reflect.String:
		fv.SetString(s)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		fv.SetInt(n)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		fv.SetUint(n)
		return nil
	default:
		return errors.New("unsupported kind in setValue")
	}
}

func (ctx *echoContext) JSON(status int, body any) error {
	return ctx.c.JSON(status, body)
}

// Result turns (data, err) into a standardized HTTP response shape.
func (ctx *echoContext) Result(data any, err error) error {
	if err == nil {
		return ctx.c.JSON(http.StatusOK, map[string]any{
			"code": errs.CodeOK,
			"data": data,
		})
	}

	var e *errs.Error
	if errors.As(err, &e) {
		status := http.StatusInternalServerError
		switch e.Code {
		case errs.CodeBadRequest:
			status = http.StatusBadRequest
		case errs.CodeNotFound:
			status = http.StatusNotFound
		}
		return ctx.c.JSON(status, map[string]any{
			"code": e.Code,
			"msg":  e.Msg,
		})
	}

	// Fallback for unknown errors.
	return ctx.c.JSON(http.StatusInternalServerError, map[string]any{
		"code": errs.CodeInternal,
		"msg":  "internal error",
	})
}
