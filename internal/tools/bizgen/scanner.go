package bizgen

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/youbuwei/doeot-go/internal/tools/shared"
)

// 扫描带注解的 endpoint。
func scanEndpoints(moduleName string) (*scanResult, error) {
	root, err := shared.FindRepoRoot()
	if err != nil {
		return nil, err
	}
	modPath, err := shared.DetectModulePath(root)
	if err != nil {
		return nil, err
	}

	endpointDir := filepath.Join(root, "internal", moduleName, "interfaces", "endpoint")
	files, err := filepath.Glob(filepath.Join(endpointDir, "*.go"))
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	var eps []endpointInfo

	for _, path := range files {
		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || len(fn.Recv.List) == 0 {
				continue
			}

			// 提取 receiver 类型名：*UserEndpoint / *OrderEndpoint
			var recvType string
			if star, ok := fn.Recv.List[0].Type.(*ast.StarExpr); ok {
				if ident, ok := star.X.(*ast.Ident); ok {
					recvType = ident.Name
				}
			}
			if recvType == "" {
				continue
			}

			info := endpointInfo{
				StructName: recvType,
				MethodName: fn.Name.Name,
			}

			if fn.Doc != nil {
				for _, c := range fn.Doc.List {
					// c.Text 形如 "// @Route  GET /users/:id"
					text := strings.TrimSpace(strings.TrimPrefix(c.Text, "//"))
					switch {
					case strings.HasPrefix(text, "@Route"):
						parts := strings.Fields(text)
						if len(parts) >= 3 {
							info.RouteMethod = parts[1]
							info.RoutePath = parts[2]
						}
					case strings.HasPrefix(text, "@RPC"):
						parts := strings.Fields(text)
						if len(parts) >= 2 {
							info.RPCMethod = parts[1]
						}
					case strings.HasPrefix(text, "@Auth"):
						parts := strings.Fields(text)
						if len(parts) >= 2 {
							info.Auth = parts[1]
						}
					case strings.HasPrefix(text, "@Tags"):
						parts := strings.Fields(text)
						if len(parts) >= 2 {
							info.Tags = parts[1:]
						}
					}
				}
			}

			// 只保存至少有一个注解的函数。
			if info.RouteMethod != "" || info.RPCMethod != "" {
				eps = append(eps, info)
			}
		}
	}

	return &scanResult{
		Endpoints: eps,
		RootDir:   root,
		ModPath:   modPath,
	}, nil
}
