# doeot-go

> 基于 Echo + GORM 的业务优先型 Go 微服务脚手架  
> 模块化、注解驱动、代码生成、HTTP + JSON-RPC、一键本地开发。

> Module: `github.com/youbuwei/doeot-go`

---

## ✨ 特性概览

- **Echo 驱动的 HTTP 服务**
    - 统一 `biz.Context` 封装请求上下文
    - 全局错误码封装（`pkg/errs`），统一返回格式
- **JSON-RPC 服务 & 多端口**
    - HTTP / RPC 分端口启动（例如 `:8080` / `:19001`）
    - 简单的 `RPCRouter` 接口抽象，支持中间件（鉴权、打点等）
- **注解 + 代码生成**
    - 在 `interfaces/endpoint` 中写业务方法 + 注解：
        - `@Route`：生成 HTTP 路由 & 请求绑定
        - `@RPC`：生成 RPC Handler
        - `@Auth` / `@Tags`：生成链路元信息（用于鉴权、监控、文档等）
    - `bizgen` 自动生成：
        - `internal/<module>/interfaces/http/zz_routes_gen.go`
        - `internal/<module>/interfaces/rpc/zz_rpc_gen.go`
- **模块化领域设计**
    - `domain` / `app` / `infra` / `interfaces` / `module`
    - `Module` 实现统一接口，支持在应用中按需注册
- **CLI 工具集合（单入口）**
    - `doeot dev` —— 本地开发（多服务 + 热更新 + HTTP 面板）
    - `doeot modgen` —— 一键生成完整业务模块骨架
    - `doeot bizgen` —— 手动根据注解生成 HTTP/RPC 包装
- **本地开发体验**
    - 监听 `internal/`、`pkg/`、`cmd/` 下的 `.go` 变更
    - 自动触发 `go generate ./...`（仅在需要时）
    - 自动重启多服务（进程组 Kill，绕过端口占用问题）
    - 内置 Dev HTTP 面板（默认 `:18080`）查看服务状态
- **ORM & MySQL**
    - 内建 GORM 集成，`infra/repo` 提供默认实现
- **配置 & .env 支持**
    - 支持 `.env` 加载，配合配置中心（如 etcd）扩展

> PS：部分特性（如配置中心、服务发现/注册）在代码中预留扩展点，可按业务节奏逐步补齐。

---

## 🧱 项目结构（核心部分）

```bash
.
├── cmd
│   ├── doeot          # 统一 CLI 入口
│   ├── dev            # dev 工具（可直接 go run 使用）
│   ├── bizgen         # 注解 -> HTTP/RPC 代码生成
│   ├── modgen         # 模块脚手架生成
│   ├── user-api       # 示例 HTTP 服务
│   └── user-rpc       # 示例 RPC 服务
├── internal
│   ├── user
│   │   ├── domain     # 领域模型 & 仓储接口
│   │   ├── app        # 应用服务（业务逻辑）
│   │   ├── infra
│   │   │   └── repo   # 基于 GORM 的仓储实现
│   │   ├── interfaces
│   │   │   └── endpoint  # 带注解的业务端点（手写）
│   │   │   └── http      # 由 bizgen 生成的 HTTP 路由
│   │   │   └── rpc       # 由 bizgen 生成的 RPC Handler
│   │   └── module.go     # Module 实现（装配依赖）
│   └── ...               # 其他业务模块（如 order）
├── pkg
│   ├── biz           # 核心上下文 & Router 封装
│   ├── boot          # 应用启动（HTTP/RPC 入口）
│   ├── errs          # 错误码体系
│   └── validate      # 请求校验封装
└── go.mod

---

## 🚀 快速开始

### 1. 克隆 & 依赖

```bash
git clone https://github.com/youbuwei/doeot-go.git
cd doeot-go

go mod tidy
```

### 2. 本地开发（dev 模式）

推荐使用统一入口 `doeot`：

```bash
# 同时跑 user-api + user-rpc，并开启 Dev HTTP 面板 :18080
go run ./cmd/doeot dev -services user-api,user-rpc -dev-http :18080
```

启动后可以看到：

```text
dev: running services: user-api, user-rpc
dev: HTTP panel: http://localhost:18080/
dev: commands: [r] restart (go generate + restart), [s] status, [q] quit
dev>
```

* 修改 `.go` 文件 → 自动热重启对应服务
* 修改 `interfaces/endpoint` / `*_endpoint.go` → 自动触发 `go generate ./...` + 重启
* 访问 Dev 面板：`http://localhost:18080/` 查看：

    * 当前注册服务
    * 每个服务状态 & PID
    * 最近一次重启 / go generate / 文件变更

---

## 🧩 生成一个新业务模块

使用聚合工具 `doeot`：

```bash
# 生成名为 order 的模块
go run ./cmd/doeot modgen -name order
```

它会自动生成：

* `internal/order/domain/domain.go`
* `internal/order/app/service.go`
* `internal/order/infra/repo/repo.go`
* `internal/order/interfaces/endpoint/order_endpoint.go`
* `internal/order/module.go`
* 调用 bizgen 生成：

    * `internal/order/interfaces/http/zz_routes_gen.go`
    * `internal/order/interfaces/rpc/zz_rpc_gen.go`

在你的服务中注册该模块，例如：

```go
import (
    "log"

    "github.com/youbuwei/doeot-go/internal/user"
    "github.com/youbuwei/doeot-go/internal/order"
    "github.com/youbuwei/doeot-go/pkg/boot"
)

func main() {
    app := boot.New("user-api")

    app.RegisterModule(user.NewModule(app.DB()))
    app.RegisterModule(order.NewModule(app.DB()))

    if err := app.Run(); err != nil {
        log.Fatal(err)
    }
}
```

---

## 📌 注解风格的 Endpoint

在 `internal/order/interfaces/endpoint/order_endpoint.go` 中：

```go
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

// OrderEndpoint 将应用服务暴露为 HTTP/RPC 端点。
type OrderEndpoint struct {
    Svc *app.OrderService
}

// @Route  GET /orders/:id
// @RPC    Order.Get
// @Auth   login
// @Desc   获取订单
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
```

然后执行（或者由 dev 自动执行）：

```bash
go run ./cmd/doeot bizgen -module order
# 或
go generate ./...
```

`bizgen` 会根据注解生成：

```go
// internal/order/interfaces/http/zz_routes_gen.go

func RegisterRoutes(r biz.Router, ep *endpoint.OrderEndpoint) {
    r.GET("/orders/:id", func(ctx biz.Context) error {
        var req endpoint.GetOrderReq
        if err := ctx.Bind(&req); err != nil {
            return ctx.Result(nil, errs.BadRequest("invalid request").WithCause(err))
        }
        if err := validate.Struct(&req); err != nil {
            return ctx.Result(nil, err)
        }
        resp, err := ep.GetOrder(ctx, &req)
        return ctx.Result(resp, err)
    }, biz.WithAuth("login"), biz.WithTags("order"), biz.WithBizTag("order.getorder"))
}
```

同理，RPC 部分也会自动生成。

---

## ✅ 统一 CLI：doeot

`cmd/doeot/main.go` 提供一个统一入口，方便记忆和使用：

```bash
# 查看 help
go run ./cmd/doeot -h

# 开发模式（dev）
go run ./cmd/doeot dev -services user-api,user-rpc -dev-http :18080

# 模块生成（modgen）
go run ./cmd/doeot modgen -name order

# 代码生成（bizgen）
go run ./cmd/doeot bizgen -module user
```

输出示例：

```text
doeot - DOEOT 项目开发工具集合

用法:
  doeot <command> [arguments]

可用命令:
  dev       启动开发模式（热更新、多服务、HTTP 面板）
  modgen    生成业务模块骨架 (domain/app/repo/endpoint/module + bizgen)
  bizgen    根据 endpoint 注解生成 HTTP/RPC 包装代码
```

---

## 🛣 Roadmap

* [ ] 配置中心集成（etcd）
* [ ] RPC 服务发现 & 注册
* [ ] 定时任务（统一调度 & 注册）
* [ ] 内置缓存封装（Redis/本地 cache）
* [ ] Swagger / OpenAPI 文档生成
* [ ] 更完善的 Auth / RBAC 组件
* [ ] 多租户 / 多环境配置管理

---

## 📄 License

MIT License.
Copyright (c) 2025.

---

## 🤝 参与贡献

欢迎 Issue / PR / 讨论：

* 新增模块模板（比如带分页、搜索条件）
* 更丰富的注解能力（幂等、幂等 Key、限流、熔断）
* Dev 面板的操作能力（Web 上一键 Restart / Bizgen / Modgen）

如果你想把自己的一套最佳实践固化到框架里，也可以直接提需求，我们可以一起把脚手架打磨成“上手就能开干业务”的形态。
