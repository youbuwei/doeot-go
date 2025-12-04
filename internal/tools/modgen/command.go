package modgen

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/youbuwei/doeot-go/pkg/cli"
)

type modgenCommand struct{}

func NewCommand() cli.Command { return &modgenCommand{} }

func (c *modgenCommand) Name() string { return "modgen" }
func (c *modgenCommand) Description() string {
	return "生成业务模块骨架 (domain/app/repo/endpoint/module + bizgen)"
}

func (c *modgenCommand) Run(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("modgen", flag.ContinueOnError)
	nameFlag := fs.String("name", "", "模块名，例如: user, order")
	fs.SetOutput(os.Stdout)

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *nameFlag == "" {
		return fmt.Errorf("请使用 -name 指定模块名，例如: doeot modgen -name order")
	}

	cfg := Config{ModuleName: *nameFlag}
	return Run(ctx, cfg)
}
