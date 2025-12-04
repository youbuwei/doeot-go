package dev

import (
	"context"
	"flag"
	"os"
	"strings"

	"github.com/youbuwei/doeot-go/pkg/cli"
)

type devCommand struct{}

func NewCommand() cli.Command { return &devCommand{} }

func (c *devCommand) Name() string { return "dev" }
func (c *devCommand) Description() string {
	return "开发模式（多服务 + 热更新 + HTTP 面板）"
}

func (c *devCommand) Run(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("dev", flag.ContinueOnError)
	servicesFlag := fs.String("services", "http-api,json-rpc", "逗号分隔服务名，如: http-api,json-rpc")
	httpAddr := fs.String("dev-http", ":18080", "HTTP 面板地址，如: :18080")
	fs.SetOutput(os.Stdout)

	if err := fs.Parse(args); err != nil {
		return err
	}

	cfg := Config{
		Services:      splitAndTrim(*servicesFlag),
		HTTPPanelAddr: *httpAddr,
	}
	return Run(ctx, cfg)
}

func splitAndTrim(s string) []string {
	raw := strings.Split(s, ",")
	var out []string
	for _, v := range raw {
		v = strings.TrimSpace(v)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}
