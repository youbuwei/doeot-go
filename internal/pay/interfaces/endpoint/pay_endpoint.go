package endpoint

//go:generate go run github.com/youbuwei/doeot-go/cmd/bizgen -module pay

import (
	"errors"

	"github.com/youbuwei/doeot-go/internal/pay/app"
	"github.com/youbuwei/doeot-go/internal/pay/domain"
	"github.com/youbuwei/doeot-go/pkg/biz"
	"github.com/youbuwei/doeot-go/pkg/errs"
)

type GetPayReq struct {
	ID int64
}

type GetPayResp struct {
	ID   int64
	Name string
}

type CreatePayReq struct {
	Name string
}

type CreatePayResp struct {
	ID   int64
	Name string
}

// PayEndpoint 将应用服务暴露为 HTTP/RPC 端点。
type PayEndpoint struct {
	Svc *app.PayService
}

// @Route  GET /pays/:id
// @RPC    Pay.Get
// @Auth   login
// @Desc   获取 pay
// @Tags   pay
func (e *PayEndpoint) GetPay(ctx biz.Context, req *GetPayReq) (*GetPayResp, error) {
	m, err := e.Svc.Get(ctx.RequestContext(), req.ID)
	if errors.Is(err, domain.ErrPayNotFound) {
		return nil, errs.NotFound("pay not found")
	}
	if err != nil {
		return nil, errs.Internal("failed to get pay").WithCause(err)
	}
	return &GetPayResp{
		ID:   m.ID,
		Name: m.Name,
	}, nil
}

// @Route  POST /pays
// @RPC    Pay.Create
// @Auth   login
// @Desc   创建 pay
// @Tags   pay
func (e *PayEndpoint) CreatePay(ctx biz.Context, req *CreatePayReq) (*CreatePayResp, error) {
	m, err := e.Svc.Create(ctx.RequestContext(), &domain.Pay{
		Name: req.Name,
	})
	if err != nil {
		return nil, errs.Internal("failed to create pay").WithCause(err)
	}
	return &CreatePayResp{
		ID:   m.ID,
		Name: m.Name,
	}, nil
}
