package main

import (
    "log"

    "github.com/youbuwei/doeot-go/internal/app"
    httpiface "github.com/youbuwei/doeot-go/internal/interfaces/http"
    "github.com/youbuwei/doeot-go/pkg/config"
    "github.com/youbuwei/doeot-go/pkg/di"
)

func main() {
    container := di.Build()

    err := container.Invoke(func(cfg *config.Config, svc app.OrderService) {
        r := httpiface.NewEngine(svc)
        if err := r.Run(":" + cfg.HTTPPort); err != nil {
            log.Fatal(err)
        }
    })
    if err != nil {
        log.Fatal(err)
    }
}
