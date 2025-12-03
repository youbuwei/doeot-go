package order

import (
	"github.com/youbuwei/doeot-go/internal/order/app"
	"github.com/youbuwei/doeot-go/internal/order/infra/repo"
	"github.com/youbuwei/doeot-go/internal/order/interfaces/endpoint"
	http "github.com/youbuwei/doeot-go/internal/order/interfaces/http"
	rpc "github.com/youbuwei/doeot-go/internal/order/interfaces/rpc"
	"github.com/youbuwei/doeot-go/pkg/biz"
	"gorm.io/gorm"
)

// Module 实现 biz.Module，用于注册 order 模块的 HTTP/RPC 路由。
type Module struct {
	ep *endpoint.OrderEndpoint
}

// NewModule 组装 order 模块的依赖关系。
func NewModule(db *gorm.DB) *Module {
	r := repo.NewRepo(db)
	svc := app.NewOrderService(r)
	ep := &endpoint.OrderEndpoint{Svc: svc}
	return &Module{ep: ep}
}

func (m *Module) Name() string { return "order" }

func (m *Module) RegisterHTTP(r biz.Router) {
	http.RegisterRoutes(r, m.ep)
}

func (m *Module) RegisterRPC(r biz.RPCRouter) {
	rpc.RegisterRPC(r, m.ep)
}
