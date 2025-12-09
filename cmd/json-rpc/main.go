package main

import (
    "log"
    "net/http"

    "github.com/youbuwei/doeot-go/internal/app"
    rpciface "github.com/youbuwei/doeot-go/internal/interfaces/rpc"
    "github.com:youbuwei/doeot-go/pkg/config"
    "github.com:youbuwei/doeot-go/pkg/di"
)

func main() {
    container := di.Build()

    err := container.Invoke(func(cfg *config.Config, svc app.OrderService) {
        handler := rpciface.NewHTTPHandler(svc)
        if err := http.ListenAndServe(":"+cfg.RPCPort, handler); err != nil {
            log.Fatal(err)
        }
    })
    if err != nil {
        log.Fatal(err)
    }
}
