package endpoint

//go:generate go run github.com/youbuwei/doeot-go/cmd/bizgen -module order

import (
	"errors"

	"github.com/youbuwei/doeot-go/internal/order/app"
	"github.com/youbuwei/doeot-go/internal/order/domain"
	"github.com/youbuwei/doeot-go/pkg/biz"
	"github.com/youbuwei/doeot-go/pkg/errs"
)

type GetOrderReq struct {
	ID int64
}

type GetOrderResp struct {
	ID   int64
	Name string
}

type CreateOrderReq struct {
	Name string
}

type CreateOrderResp struct {
	ID   int64
	Name string
}

// OrderEndpoint 将应用服务暴露为 HTTP/RPC 端点。
type OrderEndpoint struct {
	Svc *app.OrderService
}

// GetOrder
// @Route  GET /order/:id
// @RPC    Order.Get
// @Auth   login
// @Desc   获取 order
// @Tags   order
func (e *OrderEndpoint) GetOrder(ctx biz.Context, req *GetOrderReq) (*GetOrderResp, error) {
	m, err := e.Svc.Get(ctx.RequestContext(), req.ID)
	if errors.Is(err, domain.ErrOrderNotFound) {
		return nil, errs.NotFound("order not found")
	}
	if err != nil {
		return nil, errs.Internal("failed to get order").WithCause(err)
	}
	return &GetOrderResp{
		ID:   m.ID,
		Name: m.Name,
	}, nil
}

// CreateOrder
// @Route  POST /orders
// @RPC    Order.Create
// @Auth   login
// @Desc   创建 order
// @Tags   order
func (e *OrderEndpoint) CreateOrder(ctx biz.Context, req *CreateOrderReq) (*CreateOrderResp, error) {
	m, err := e.Svc.Create(ctx.RequestContext(), &domain.Order{
		Name: req.Name,
	})
	if err != nil {
		return nil, errs.Internal("failed to create order").WithCause(err)
	}
	return &CreateOrderResp{
		ID:   m.ID,
		Name: m.Name,
	}, nil
}
