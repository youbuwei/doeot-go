package config

import "os"

// Config 基础配置
type Config struct {
	HTTP_PORT string
	MYSQL_DSN string
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// Load 从环境变量加载配置（带默认值）
func Load() (*Config, error) {
	cfg := &Config{
		HTTP_PORT: getEnv("HTTP_PORT", "8080"),
		MYSQL_DSN: getEnv("MYSQL_DSN", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"),
	}
	return cfg, nil
}
