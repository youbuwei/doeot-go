package bizgen

// 从 endpoint 源码扫描出来的基础信息。
type endpointInfo struct {
	StructName  string   // 如 UserEndpoint / OrderEndpoint
	MethodName  string   // 如 GetUser
	RouteMethod string   // "GET"/"POST"/...
	RoutePath   string   // "/users/:id"
	RPCMethod   string   // "User.Get"
	Auth        string   // 来自 @Auth
	Tags        []string // 来自 @Tags
}

// 扫描结果：包含模块根目录等信息。
type scanResult struct {
	Endpoints []endpointInfo
	RootDir   string // 仓库根目录（包含 go.mod）
	ModPath   string // go.mod 里的 module 路径
}

// 模板使用的结构（HTTP）。
type httpEndpointData struct {
	MethodName string
	HTTPMethod string
	RoutePath  string
	Options    string // "biz.WithAuth(...), biz.WithTags(...), biz.WithBizTag(...)"
}

type httpTemplateData struct {
	ModPath      string
	Module       string
	EndpointType string
	Endpoints    []httpEndpointData
}

// 模板使用的结构（RPC）。
type rpcEndpointData struct {
	MethodName string
	RPCMethod  string
	Options    string
}

type rpcTemplateData struct {
	ModPath      string
	Module       string
	EndpointType string
	Endpoints    []rpcEndpointData
}
