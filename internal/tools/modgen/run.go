package modgen

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/youbuwei/doeot-go/internal/tools/bizgen"
	"github.com/youbuwei/doeot-go/internal/tools/shared"
)

// Config 在 config.go 里已定义，这里直接使用。
// type Config struct{ ModuleName string }

// templateData 用于 domain/app/repo/endpoint/module 模板。
type templateData struct {
	ModPath    string // go.mod 里的 module 路径，例如 github.com/youbuwei/doeot-go
	ModuleName string // 模块名，形如 "user"、"order"、"pay"
	TypeName   string // 导出的类型名，形如 "User"、"Order"、"Pay"
}

// modulesTemplateData 用于 modules.tmpl。
type modulesTemplateData struct {
	ModPath string   // module 路径
	Modules []string // 所有业务模块名，例如 ["user", "order", "pay"]
}

// Run 根据 Config 生成模块骨架，并调用 bizgen + 模块注册表生成。
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

	data := templateData{
		ModPath:    modPath,
		ModuleName: name,
		TypeName:   typeName,
	}

	// 1. 领域层（只在文件不存在时创建）
	if err := genFromTemplate(root, domainTmpl,
		filepath.Join("internal", name, "domain", "domain.go"), data, false); err != nil {
		return err
	}

	// 2. 应用层
	if err := genFromTemplate(root, appTmpl,
		filepath.Join("internal", name, "app", "service.go"), data, false); err != nil {
		return err
	}

	// 3. 仓储实现
	if err := genFromTemplate(root, repoTmpl,
		filepath.Join("internal", name, "infra", "repo", "repo.go"), data, false); err != nil {
		return err
	}

	// 4. endpoint（带注解）
	if err := genFromTemplate(root, endpointTmpl,
		filepath.Join("internal", name, "interfaces", "endpoint", name+"_endpoint.go"), data, false); err != nil {
		return err
	}

	// 5. module 入口
	if err := genFromTemplate(root, moduleTmpl,
		filepath.Join("internal", name, "module.go"), data, false); err != nil {
		return err
	}

	// 6. 调用 bizgen 生成 HTTP/RPC wrapper
	if err := bizgen.Run(ctx, bizgen.Config{ModuleName: name}); err != nil {
		fmt.Fprintf(os.Stderr, "modgen: internal bizgen failed, fallback to go run cmd/bizgen: %v\n", err)
		cmd := exec.Command("go", "run", filepath.Join(root, "cmd", "bizgen"), "-module", name)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err2 := cmd.Run(); err2 != nil {
			return fmt.Errorf("fallback bizgen failed: %w", err2)
		}
	}

	// 7. 生成模块注册表 internal/modules/zz_modules_gen.go（总是允许覆盖）
	if err := generateModulesRegistry(root, modPath); err != nil {
		return err
	}

	return nil
}

// 把 s_foo / foo_bar 转成 CamelCase：sFoo / FooBar。
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

// 渲染模板并写入文件：
//   - 如果文件不存在：创建，并输出 "modgen: created <path>"
//   - 如果存在且 overwrite=false：不改，不输出
//   - 如果存在且 overwrite=true：
//   - 内容不变：不输出
//   - 内容变化：覆盖，并输出 "modgen: updated <path>"
func genFromTemplate(root string, tmpl *template.Template, relPath string, data any, overwrite bool) error {
	full := filepath.Join(root, relPath)

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("execute template for %s: %w", relPath, err)
	}
	newContent := buf.Bytes()

	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", relPath, err)
	}

	_, err := os.Stat(full)
	if errors.Is(err, os.ErrNotExist) {
		// 首次生成
		if err := os.WriteFile(full, newContent, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", relPath, err)
		}
		fmt.Printf("modgen: created %s\n", relPath)
		return nil
	}
	if err != nil {
		return fmt.Errorf("stat %s: %w", relPath, err)
	}

	// 文件已存在
	if !overwrite {
		// skeleton 文件走“只生成一次”的策略，不再改动
		return nil
	}

	oldContent, err := os.ReadFile(full)
	if err != nil {
		return fmt.Errorf("read %s: %w", relPath, err)
	}
	if bytes.Equal(oldContent, newContent) {
		// 内容未变化，无需输出
		return nil
	}

	if err := os.WriteFile(full, newContent, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", relPath, err)
	}
	fmt.Printf("modgen: updated %s\n", relPath)
	return nil
}

// 扫描 internal/*/module.go 生成 internal/modules/zz_modules_gen.go。
func generateModulesRegistry(root, modPath string) error {
	pattern := filepath.Join(root, "internal", "*", "module.go")
	paths, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	if len(paths) == 0 {
		// 还没有任何模块，跳过。
		return nil
	}

	modSet := make(map[string]struct{})
	for _, p := range paths {
		dir := filepath.Dir(p)
		mod := filepath.Base(dir)
		// 排除 internal/modules、internal/tools 等非业务目录。
		if mod == "modules" || mod == "tools" {
			continue
		}
		modSet[mod] = struct{}{}
	}

	if len(modSet) == 0 {
		return nil
	}

	var mods []string
	for m := range modSet {
		mods = append(mods, m)
	}
	sort.Strings(mods)

	data := modulesTemplateData{
		ModPath: modPath,
		Modules: mods,
	}

	relPath := filepath.Join("internal", "modules", "zz_modules_gen.go")
	// 注册表允许覆盖，因为模块列表会随时间变化。
	return genFromTemplate(root, modulesTmpl, relPath, data, true)
}
