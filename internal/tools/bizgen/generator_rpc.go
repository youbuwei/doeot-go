package bizgen

import (
	"os"
	"path/filepath"
	"strings"
)

// 生成 RPC 包装代码。
func generateRPC(res *scanResult, module string) error {
	endpointType := "Endpoint"
	if len(res.Endpoints) > 0 && res.Endpoints[0].StructName != "" {
		endpointType = res.Endpoints[0].StructName
	}

	var eps []rpcEndpointData
	for _, e := range res.Endpoints {
		if e.RPCMethod == "" {
			continue
		}
		bizTag := strings.ToLower(module + "." + strings.ToLower(e.MethodName))
		opts := buildOptions(e.Auth, e.Tags, bizTag)

		eps = append(eps, rpcEndpointData{
			MethodName: e.MethodName,
			RPCMethod:  e.RPCMethod,
			Options:    opts,
		})
	}

	if len(eps) == 0 {
		return nil
	}

	data := rpcTemplateData{
		ModPath:      res.ModPath,
		Module:       module,
		EndpointType: endpointType,
		Endpoints:    eps,
	}

	rpcDir := filepath.Join(res.RootDir, "internal", module, "interfaces", "rpc")
	if err := os.MkdirAll(rpcDir, 0o755); err != nil {
		return err
	}
	path := filepath.Join(rpcDir, "zz_rpc_gen.go")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return rpcTmpl.Execute(f, data)
}
