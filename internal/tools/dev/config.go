package dev

// Config 是 dev 工具的配置。
type Config struct {
	Services      []string // 需要托管的服务，例如 http-api, json-rpc
	HTTPPanelAddr string   // HTTP 面板地址，例如 :18080
}
