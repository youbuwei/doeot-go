package bizgen

import (
	"context"
	"flag"
	"os"

	"github.com/youbuwei/doeot-go/pkg/cli"
)

type bizgenCommand struct{}

func NewCommand() cli.Command { return &bizgenCommand{} }

func (c *bizgenCommand) Name() string { return "bizgen" }
func (c *bizgenCommand) Description() string {
	return "根据 endpoint 注解生成 HTTP/RPC 包装代码"
}

func (c *bizgenCommand) Run(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("bizgen", flag.ContinueOnError)
	module := fs.String("module", "user", "业务模块名，例如: user, order")
	fs.SetOutput(os.Stdout)

	if err := fs.Parse(args); err != nil {
		return err
	}
	cfg := Config{ModuleName: *module}
	return Run(ctx, cfg)
}
