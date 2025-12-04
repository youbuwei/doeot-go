package bizgen

type endpointInfo struct {
	StructName  string
	MethodName  string
	RouteMethod string
	RoutePath   string
	RPCMethod   string
	Auth        string
	Tags        []string
}

type scanResult struct {
	Endpoints []endpointInfo
	RootDir   string
	ModPath   string
}
