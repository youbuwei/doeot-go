package main

import (
	"context"
	"os"

	"github.com/youbuwei/doeot-go/internal/tools/bizgen"
	"github.com/youbuwei/doeot-go/internal/tools/dev"
	"github.com/youbuwei/doeot-go/internal/tools/modgen"
	"github.com/youbuwei/doeot-go/pkg/cli"
)

func main() {
	app := cli.NewApp("doeot", "DOEOT 项目开发工具集合")

	// 内置工具命令。
	app.Register(dev.NewCommand())
	app.Register(modgen.NewCommand())
	app.Register(bizgen.NewCommand())

	// 将来这里还可以注册业务模块的命令:
	// app.RegisterProvider(usercmd.NewUserCommands())

	os.Exit(app.Run(context.Background(), os.Args[1:]))
}
