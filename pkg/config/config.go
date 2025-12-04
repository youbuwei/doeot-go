package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// MySQLConfig holds database settings.
type MySQLConfig struct {
	DSN        string
	MaxIdle    int
	MaxOpen    int
	MaxLifeMin int
}

// HTTPConfig holds HTTP server settings.
type HTTPConfig struct {
	Addr string
}

// RPCConfig holds RPC server settings.
type RPCConfig struct {
	Addr string
}

// AppConfig groups all configuration parts.
type AppConfig struct {
	MySQL MySQLConfig
	HTTP  HTTPConfig
	RPC   RPCConfig
}

// Load returns config for a given service name.
// It loads .env (if present), then reads from environment variables with sane defaults.
// Service-specific overrides can be done via env if needed.
func Load(serviceName string) AppConfig {
	// Best-effort load .env; ignore error if file is missing.
	_ = godotenv.Load()

	mysqlDSN := getenv("MYSQL_DSN",
		"root:root@tcp(127.0.0.1:3306)/mall?parseTime=true&loc=Local",
	)
	mysqlMaxIdle := getenvInt("MYSQL_MAX_IDLE", 10)
	mysqlMaxOpen := getenvInt("MYSQL_MAX_OPEN", 50)
	mysqlMaxLifeMin := getenvInt("MYSQL_MAX_LIFE_MIN", 60)

	httpAddr := os.Getenv("HTTP_ADDR")
	rpcAddr := os.Getenv("RPC_ADDR")

	// Some sane defaults based on serviceName if env not set explicitly.
	switch serviceName {
	case "http-api":
		if httpAddr == "" {
			httpAddr = ":8080"
		}
		// http-api is HTTP-only by default
		rpcAddr = ""
	case "json-rpc":
		if rpcAddr == "" {
			rpcAddr = ":19001"
		}
		// json-rpc is RPC-only by default
		httpAddr = ""
	}

	return AppConfig{
		MySQL: MySQLConfig{
			DSN:        mysqlDSN,
			MaxIdle:    mysqlMaxIdle,
			MaxOpen:    mysqlMaxOpen,
			MaxLifeMin: mysqlMaxLifeMin,
		},
		HTTP: HTTPConfig{
			Addr: httpAddr,
		},
		RPC: RPCConfig{
			Addr: rpcAddr,
		},
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Printf("invalid int env %s=%s, use default %d", key, v, def)
			return def
		}
		return n
	}
	return def
}
