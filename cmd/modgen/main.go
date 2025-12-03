package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// go run ./cmd/modgen -name order
func main() {
	nameFlag := flag.String("name", "", "模块名，例如: user, order")
	flag.Parse()

	if *nameFlag == "" {
		log.Fatal("请使用 -name 指定模块名，例如: go run ./cmd/modgen -name order")
	}

	moduleName := strings.ToLower(*nameFlag) // internal/order
	typeName := toCamel(moduleName)          // Order

	modulePath, err := detectModulePath()
	if err != nil {
		log.Fatalf("解析 go.mod 失败: %v", err)
	}

	log.Printf("modgen: module=%s, type=%s, modulePath=%s", moduleName, typeName, modulePath)

	if err := generateDomain(moduleName, typeName); err != nil {
		log.Fatalf("生成 domain 失败: %v", err)
	}
	if err := generateApp(modulePath, moduleName, typeName); err != nil {
		log.Fatalf("生成 app 失败: %v", err)
	}
	if err := generateRepo(modulePath, moduleName, typeName); err != nil {
		log.Fatalf("生成 repo 失败: %v", err)
	}
	if err := generateEndpoint(modulePath, moduleName, typeName); err != nil {
		log.Fatalf("生成 endpoint 失败: %v", err)
	}
	if err := generateModule(modulePath, moduleName, typeName); err != nil {
		log.Fatalf("生成 module.go 失败: %v", err)
	}

	if err := runBizgen(modulePath, moduleName); err != nil {
		log.Printf("modgen: 运行 bizgen 失败（你也可以手动执行: go run %s/cmd/bizgen -module %s）：%v",
			modulePath, moduleName, err)
	} else {
		log.Printf("modgen: 已执行 bizgen 生成 HTTP/RPC wrapper")
	}

	log.Println("modgen: done.")
}

// 从 go.mod 中解析 module 路径
func detectModulePath() (string, error) {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}
	return "", fmt.Errorf("未在 go.mod 中找到 module 声明")
}

// 简单 camel 转换
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

// 写文件：若文件已存在则跳过
func writeFileOnce(path string, content string) error {
	if _, err := os.Stat(path); err == nil {
		log.Printf("modgen: 跳过已存在文件 %s", path)
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return err
	}
	log.Printf("modgen: 创建文件 %s", path)
	return nil
}

// --- 各层生成 ---

func generateDomain(moduleName, typeName string) error {
	path := filepath.Join("internal", moduleName, "domain", "domain.go")
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
	return writeFileOnce(path, content)
}

func generateApp(modulePath, moduleName, typeName string) error {
	path := filepath.Join("internal", moduleName, "app", "service.go")
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
`, modulePath, moduleName, typeName)
	return writeFileOnce(path, content)
}

func generateRepo(modulePath, moduleName, typeName string) error {
	path := filepath.Join("internal", moduleName, "infra", "repo", "repo.go")
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
`, modulePath, moduleName, typeName)
	return writeFileOnce(path, content)
}

func generateEndpoint(modulePath, moduleName, typeName string) error {
	path := filepath.Join("internal", moduleName, "interfaces", "endpoint", moduleName+"_endpoint.go")

	header := fmt.Sprintf(`package endpoint

import (
	"errors"

	"%[1]s/internal/%[2]s/app"
	"%[1]s/internal/%[2]s/domain"
	"%[1]s/pkg/biz"
	"%[1]s/pkg/errs"
)
`, modulePath, moduleName)

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

	return writeFileOnce(path, header+types+methods)
}

func generateModule(modulePath, moduleName, typeName string) error {
	path := filepath.Join("internal", moduleName, "module.go")
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
`, modulePath, moduleName, typeName)
	return writeFileOnce(path, content)
}

func runBizgen(modulePath, moduleName string) error {
	// 使用 module import path，而不是相对 ./cmd/bizgen
	cmd := exec.Command("go", "run", modulePath+"/cmd/bizgen", "-module", moduleName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
