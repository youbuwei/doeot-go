package modgen

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/youbuwei/doeot-go/internal/tools/bizgen"
	"github.com/youbuwei/doeot-go/internal/tools/shared"
)

// Run 根据 Config 生成模块骨架，并调用 bizgen 生成 HTTP/RPC 包装代码。
func Run(ctx context.Context, cfg Config) error {
	name := strings.ToLower(strings.TrimSpace(cfg.ModuleName))
	if name == "" {
		return fmt.Errorf("module name is empty")
	}
	typeName := toCamel(name)

	root, err := shared.FindRepoRoot()
	if err != nil {
		return err
	}
	modPath, err := shared.DetectModulePath(root)
	if err != nil {
		return err
	}

	if err := generateDomain(root, modPath, name, typeName); err != nil {
		return err
	}
	if err := generateApp(root, modPath, name, typeName); err != nil {
		return err
	}
	if err := generateRepo(root, modPath, name, typeName); err != nil {
		return err
	}
	if err := generateEndpoint(root, modPath, name, typeName); err != nil {
		return err
	}
	if err := generateModule(root, modPath, name, typeName); err != nil {
		return err
	}

	// 优先使用内部库调用 bizgen（更快）；如果失败，退回 go run cmd/bizgen。
	if err := bizgen.Run(ctx, bizgen.Config{ModuleName: name}); err != nil {
		fmt.Fprintf(os.Stderr, "modgen: internal bizgen failed, fallback to go run cmd/bizgen: %v\n", err)
		cmd := exec.Command("go", "run", filepath.Join(root, "cmd", "bizgen"), "-module", name)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err2 := cmd.Run(); err2 != nil {
			return fmt.Errorf("fallback bizgen failed: %w", err2)
		}
	}

	return nil
}

func toCamel(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + strings.ToLower(p[1:])
	}
	return strings.Join(parts, "")
}

func generateDomain(root, modPath, moduleName, typeName string) error {
	path := filepath.Join(root, "internal", moduleName, "domain", "domain.go")
	content := fmt.Sprintf(`package domain

import "context"

// %[2]s 是 %[1]s 模块的领域模型示例，你可以按需扩展字段。
type %[2]s struct {
    ID   int64
    Name string
}

// NotFoundError 用于标识未找到该资源。
type NotFoundError struct {
    msg string
}

func (e *NotFoundError) Error() string { return e.msg }

// Err%[2]sNotFound 在仓储查不到数据时返回。
var Err%[2]sNotFound = &NotFoundError{msg: "%[1]s not found"}

// Repo 抽象了针对 %[2]s 的持久化操作。
type Repo interface {
    FindByID(ctx context.Context, id int64) (*%[2]s, error)
    Create(ctx context.Context, m *%[2]s) (*%[2]s, error)
    List(ctx context.Context) ([]*%[2]s, error)
}
`, moduleName, typeName)
	return shared.WriteFileOnce(path, []byte(content))
}

func generateApp(root, modPath, moduleName, typeName string) error {
	path := filepath.Join(root, "internal", moduleName, "app", "service.go")
	content := fmt.Sprintf(`package app

import (
    "context"

    "%[1]s/internal/%[2]s/domain"
)

// %[3]sService 封装了围绕 %[3]s 的业务逻辑。
type %[3]sService struct {
    repo domain.Repo
}

func New%[3]sService(repo domain.Repo) *%[3]sService {
    return &%[3]sService{repo: repo}
}

func (s *%[3]sService) Get(ctx context.Context, id int64) (*domain.%[3]s, error) {
    return s.repo.FindByID(ctx, id)
}

func (s *%[3]sService) Create(ctx context.Context, m *domain.%[3]s) (*domain.%[3]s, error) {
    return s.repo.Create(ctx, m)
}

func (s *%[3]sService) List(ctx context.Context) ([]*domain.%[3]s, error) {
    return s.repo.List(ctx)
}
`, modPath, moduleName, typeName)
	return shared.WriteFileOnce(path, []byte(content))
}

func generateRepo(root, modPath, moduleName, typeName string) error {
	path := filepath.Join(root, "internal", moduleName, "infra", "repo", "repo.go")
	content := fmt.Sprintf(`package repo

import (
    "context"
    "errors"

    "%[1]s/internal/%[2]s/domain"
    "gorm.io/gorm"
)

// %[3]sModel 是 %[2]s 模块的 GORM 模型。
type %[3]sModel struct {
    ID   int64
    Name string
}

func (%[3]sModel) TableName() string { return "%[2]ss" }

// Repo 是基于 GORM 的 domain.Repo 实现。
type Repo struct {
    db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
    return &Repo{db: db}
}

func (r *Repo) FindByID(ctx context.Context, id int64) (*domain.%[3]s, error) {
    var m %[3]sModel
    if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, domain.Err%[3]sNotFound
        }
        return nil, err
    }
    return &domain.%[3]s{
        ID:   m.ID,
        Name: m.Name,
    }, nil
}

func (r *Repo) Create(ctx context.Context, d *domain.%[3]s) (*domain.%[3]s, error) {
    m := %[3]sModel{
        ID:   d.ID,
        Name: d.Name,
    }
    if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
        return nil, err
    }
    return &domain.%[3]s{
        ID:   m.ID,
        Name: m.Name,
    }, nil
}

func (r *Repo) List(ctx context.Context) ([]*domain.%[3]s, error) {
    var rows []%[3]sModel
    if err := r.db.WithContext(ctx).Find(&rows).Error; err != nil {
        return nil, err
    }
    res := make([]*domain.%[3]s, 0, len(rows))
    for _, m := range rows {
        res = append(res, &domain.%[3]s{
            ID:   m.ID,
            Name: m.Name,
        })
    }
    return res, nil
}
`, modPath, moduleName, typeName)
	return shared.WriteFileOnce(path, []byte(content))
}

func generateEndpoint(root, modPath, moduleName, typeName string) error {
	path := filepath.Join(root, "internal", moduleName, "interfaces", "endpoint", moduleName+"_endpoint.go")
	header := fmt.Sprintf(`package endpoint

//go:generate go run github.com/youbuwei/doeot-go/cmd/bizgen -module %[2]s

import (
    "errors"

    "%[1]s/internal/%[2]s/app"
    "%[1]s/internal/%[2]s/domain"
    "%[1]s/pkg/biz"
    "%[1]s/pkg/errs"
)
`, modPath, moduleName)

	types := fmt.Sprintf(`
type Get%[1]sReq struct {
    ID int64
}

type Get%[1]sResp struct {
    ID   int64
    Name string
}

type Create%[1]sReq struct {
    Name string
}

type Create%[1]sResp struct {
    ID   int64
    Name string
}

// %[1]sEndpoint 将应用服务暴露为 HTTP/RPC 端点。
type %[1]sEndpoint struct {
    Svc *app.%[1]sService
}
`, typeName)

	methods := fmt.Sprintf(`
// @Route  GET /%[2]ss/:id
// @RPC    %[1]s.Get
// @Auth   login
// @Desc   获取 %[2]s
// @Tags   %[2]s
func (e *%[1]sEndpoint) Get%[1]s(ctx biz.Context, req *Get%[1]sReq) (*Get%[1]sResp, error) {
    m, err := e.Svc.Get(ctx.RequestContext(), req.ID)
    if errors.Is(err, domain.Err%[1]sNotFound) {
        return nil, errs.NotFound("%[2]s not found")
    }
    if err != nil {
        return nil, errs.Internal("failed to get %[2]s").WithCause(err)
    }
    return &Get%[1]sResp{
        ID:   m.ID,
        Name: m.Name,
    }, nil
}

// @Route  POST /%[2]ss
// @RPC    %[1]s.Create
// @Auth   login
// @Desc   创建 %[2]s
// @Tags   %[2]s
func (e *%[1]sEndpoint) Create%[1]s(ctx biz.Context, req *Create%[1]sReq) (*Create%[1]sResp, error) {
    m, err := e.Svc.Create(ctx.RequestContext(), &domain.%[1]s{
        Name: req.Name,
    })
    if err != nil {
        return nil, errs.Internal("failed to create %[2]s").WithCause(err)
    }
    return &Create%[1]sResp{
        ID:   m.ID,
        Name: m.Name,
    }, nil
}
`, typeName, moduleName)

	return shared.WriteFileOnce(path, []byte(header+types+methods))
}

func generateModule(root, modPath, moduleName, typeName string) error {
	path := filepath.Join(root, "internal", moduleName, "module.go")
	content := fmt.Sprintf(`package %[2]s

import (
    "%[1]s/internal/%[2]s/app"
    "%[1]s/internal/%[2]s/infra/repo"
    "%[1]s/internal/%[2]s/interfaces/endpoint"
    http "%[1]s/internal/%[2]s/interfaces/http"
    rpc "%[1]s/internal/%[2]s/interfaces/rpc"
    "%[1]s/pkg/biz"
    "gorm.io/gorm"
)

// Module 实现 biz.Module，用于注册 %[2]s 模块的 HTTP/RPC 路由。
type Module struct {
    ep *endpoint.%[3]sEndpoint
}

// NewModule 组装 %[2]s 模块的依赖关系。
func NewModule(db *gorm.DB) *Module {
    r := repo.NewRepo(db)
    svc := app.New%[3]sService(r)
    ep := &endpoint.%[3]sEndpoint{Svc: svc}
    return &Module{ep: ep}
}

func (m *Module) Name() string { return "%[2]s" }

func (m *Module) RegisterHTTP(r biz.Router) {
    http.RegisterRoutes(r, m.ep)
}

func (m *Module) RegisterRPC(r biz.RPCRouter) {
    rpc.RegisterRPC(r, m.ep)
}
`, modPath, moduleName, typeName)
	return shared.WriteFileOnce(path, []byte(content))
}
