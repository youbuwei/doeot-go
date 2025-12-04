package cli

import (
	"context"
	"fmt"
	"os"
	"sort"
)

// App 是一个简单的 CLI 应用，用于注册命令并分发执行。
type App struct {
	name        string
	description string
	cmds        map[string]Command
}

func NewApp(name, description string) *App {
	return &App{
		name:        name,
		description: description,
		cmds:        make(map[string]Command),
	}
}

// Register 注册一个命令，如果同名会 panic，避免冲突。
func (a *App) Register(cmd Command) {
	name := cmd.Name()
	if name == "" {
		panic("cli: command name empty")
	}
	if _, ok := a.cmds[name]; ok {
		panic("cli: duplicate command: " + name)
	}
	a.cmds[name] = cmd
}

// RegisterProvider 让模块一次性注册多条命令。
func (a *App) RegisterProvider(p CommandProvider) {
	for _, cmd := range p.Commands() {
		a.Register(cmd)
	}
}

// Run 根据 args[0] 选择执行命令。
// args 不包括程序名本身，例如 os.Args[1:].
func (a *App) Run(ctx context.Context, args []string) int {
	if len(args) == 0 {
		a.printUsage()
		return 0
	}

	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		a.printUsage()
		return 0
	}

	name := args[0]
	cmd, ok := a.cmds[name]
	if !ok {
		fmt.Fprintf(os.Stderr, "%s: unknown command %q\n\n", a.name, name)
		a.printUsage()
		return 1
	}

	if err := cmd.Run(ctx, args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%s %s: %v\n", a.name, name, err)
		return 1
	}
	return 0
}

func (a *App) printUsage() {
	fmt.Fprintf(os.Stderr, "%s - %s\n\n", a.name, a.description)
	fmt.Fprintf(os.Stderr, "用法:\n  %s <command> [arguments]\n\n", a.name)
	fmt.Fprintf(os.Stderr, "可用命令:\n")

	names := make([]string, 0, len(a.cmds))
	for name := range a.cmds {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		fmt.Fprintf(os.Stderr, "  %-16s %s\n", name, a.cmds[name].Description())
	}

	fmt.Fprintln(os.Stderr, `
示例:
  doeot dev -services http-api,json-rpc -dev-http :18080
  doeot modgen -name order
  doeot bizgen -module user

提示:
  每个子命令通常也支持 -h/--help 查看自己的参数。`)
}
