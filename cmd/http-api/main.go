package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/youbuwei/doeot-go/internal/di"
	"github.com/youbuwei/doeot-go/pkg/config"
)

func main() {
	container := di.Build()

	err := container.Invoke(func(cfg *config.Config, engine *gin.Engine) {
		if err := engine.Run(":" + cfg.HTTP_PORT); err != nil {
			log.Fatal(err)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
}
