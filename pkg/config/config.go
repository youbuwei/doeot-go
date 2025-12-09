package config

import "os"

// Config 基础配置
type Config struct {
    HTTPPort string
    RPCPort  string
    DBPath   string
}

// Load 从环境变量加载配置（带默认值）
func Load() (*Config, error) {
    cfg := &Config{
        HTTPPort: "8080",
        RPCPort:  "8090",
        DBPath:   "data.db",
    }

    if v := os.Getenv("HTTP_PORT"); v != "" {
        cfg.HTTPPort = v
    }
    if v := os.Getenv("RPC_PORT"); v != "" {
        cfg.RPCPort = v
    }
    if v := os.Getenv("DB_PATH"); v != "" {
        cfg.DBPath = v
    }

    return cfg, nil
}
