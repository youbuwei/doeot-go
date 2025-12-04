package bizgen

import (
	"context"
	"log"
)

// Run 是 bizgen 的主入口，由命令行或 go:generate 调用。
func Run(ctx context.Context, cfg Config) error {
	res, err := scanEndpoints(cfg.ModuleName)
	if err != nil {
		return err
	}
	if len(res.Endpoints) == 0 {
		log.Printf("bizgen: no annotated endpoints found for module %s", cfg.ModuleName)
		return nil
	}

	if err := generateHTTP(res, cfg.ModuleName); err != nil {
		return err
	}
	if err := generateRPC(res, cfg.ModuleName); err != nil {
		return err
	}
	return nil
}
