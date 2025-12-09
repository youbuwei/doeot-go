// cmd/http-api/main.go
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/youbuwei/doeot-go/pkg/boot"
	"github.com/youbuwei/doeot-go/pkg/config"
)

func main() {
	// 加载 .env（如果没文件会忽略错误）
	_ = godotenv.Load(".env")

	container := boot.BuildContainer()

	err := container.Invoke(func(cfg *config.Config, engine *gin.Engine) {
		if err := engine.Run(":" + cfg.HTTPPort); err != nil {
			log.Fatal(err)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
}
