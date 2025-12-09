package config

import "os"

// Config 基础配置
type Config struct {
	HTTPPort string `env:"HTTP_PORT"`
	MySQLDSN string `env:"MYSQL_DSN"`
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
		HTTPPort: getEnv("HTTP_PORT", "8080"),
		MySQLDSN: getEnv("MYSQL_DSN", "root:123456@tcp(127.0.0.1:3306)/doeot?charset=utf8mb4&parseTime=True&loc=Local"),
	}
	return cfg, nil
}
