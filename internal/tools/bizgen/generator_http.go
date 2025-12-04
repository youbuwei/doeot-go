package bizgen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 生成 HTTP 路由包装代码。
func generateHTTP(res *scanResult, module string) error {
	// 先把需要的数据准备好，如果没有 HTTP 端点，就直接返回。
	endpointType := "Endpoint"
	if len(res.Endpoints) > 0 && res.Endpoints[0].StructName != "" {
		endpointType = res.Endpoints[0].StructName
	}

	var eps []httpEndpointData
	for _, e := range res.Endpoints {
		if e.RouteMethod == "" || e.RoutePath == "" {
			continue
		}
		method := strings.ToUpper(e.RouteMethod)
		switch method {
		case "GET", "POST", "PUT", "DELETE":
		default:
			continue
		}

		bizTag := strings.ToLower(module + "." + strings.ToLower(e.MethodName))
		opts := buildOptions(e.Auth, e.Tags, bizTag)

		eps = append(eps, httpEndpointData{
			MethodName: e.MethodName,
			HTTPMethod: method,
			RoutePath:  e.RoutePath,
			Options:    opts,
		})
	}

	if len(eps) == 0 {
		return nil
	}

	data := httpTemplateData{
		ModPath:      res.ModPath,
		Module:       module,
		EndpointType: endpointType,
		Endpoints:    eps,
	}

	httpDir := filepath.Join(res.RootDir, "internal", module, "interfaces", "http")
	if err := os.MkdirAll(httpDir, 0o755); err != nil {
		return err
	}
	path := filepath.Join(httpDir, "zz_routes_gen.go")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return httpTmpl.Execute(f, data)
}

// 构造中间件 Options 字符串，例如：
// "biz.WithAuth("login"), biz.WithTags("user"), biz.WithBizTag("user.getuser")"
func buildOptions(auth string, tags []string, bizTag string) string {
	var opts []string
	if auth != "" {
		opts = append(opts, fmt.Sprintf("biz.WithAuth(%q)", auth))
	}
	if len(tags) > 0 {
		qs := make([]string, 0, len(tags))
		for _, t := range tags {
			qs = append(qs, fmt.Sprintf("%q", t))
		}
		opts = append(opts, fmt.Sprintf("biz.WithTags(%s)", strings.Join(qs, ", ")))
	}
	opts = append(opts, fmt.Sprintf("biz.WithBizTag(%q)", bizTag))
	return strings.Join(opts, ", ")
}
