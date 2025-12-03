package boot

import (
    "github.com/youbuwei/doeot-go/pkg/biz"
    "github.com/youbuwei/doeot-go/pkg/config"
    "github.com/youbuwei/doeot-go/pkg/orm"
    "gorm.io/gorm"
)

// App holds shared infrastructure objects and registered modules.
type App struct {
    name    string
    cfg     config.AppConfig
    db      *gorm.DB
    modules []biz.Module
}

// New creates a new application for the given service name.
func New(serviceName string) *App {
    cfg := config.Load(serviceName)
    db := orm.NewMySQL(cfg.MySQL)

    return &App{
        name: serviceName,
        cfg:  cfg,
        db:   db,
    }
}

// DB exposes the shared *gorm.DB instance to wiring code in main or modules.
func (a *App) DB() *gorm.DB {
    return a.db
}

// RegisterModule registers a business module which can attach HTTP/RPC routes.
func (a *App) RegisterModule(m biz.Module) {
    a.modules = append(a.modules, m)
}

// Run starts whichever transports are configured (HTTP or RPC).
func (a *App) Run() error {
    switch {
    case a.cfg.HTTP.Addr != "" && a.cfg.RPC.Addr == "":
        return a.runHTTP()
    case a.cfg.HTTP.Addr == "" && a.cfg.RPC.Addr != "":
        return a.runRPC()
    case a.cfg.HTTP.Addr != "" && a.cfg.RPC.Addr != "":
        // For demo we just start HTTP; in real world you'd start both in goroutines.
        return a.runHTTP()
    default:
        return nil
    }
}
