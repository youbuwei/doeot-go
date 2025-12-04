package pay

import (
	"github.com/youbuwei/doeot-go/internal/pay/app"
	"github.com/youbuwei/doeot-go/internal/pay/infra/repo"
	"github.com/youbuwei/doeot-go/internal/pay/interfaces/endpoint"
	http "github.com/youbuwei/doeot-go/internal/pay/interfaces/http"
	rpc "github.com/youbuwei/doeot-go/internal/pay/interfaces/rpc"
	"github.com/youbuwei/doeot-go/pkg/biz"
	"gorm.io/gorm"
)

// Module 实现 biz.Module，用于注册 pay 模块的 HTTP/RPC 路由。
type Module struct {
	ep *endpoint.PayEndpoint
}

// NewModule 组装 pay 模块的依赖关系。
func NewModule(db *gorm.DB) *Module {
	r := repo.NewRepo(db)
	svc := app.NewPayService(r)
	ep := &endpoint.PayEndpoint{Svc: svc}
	return &Module{ep: ep}
}

func (m *Module) Name() string { return "pay" }

func (m *Module) RegisterHTTP(r biz.Router) {
	http.RegisterRoutes(r, m.ep)
}

func (m *Module) RegisterRPC(r biz.RPCRouter) {
	rpc.RegisterRPC(r, m.ep)
}
