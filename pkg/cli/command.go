package cli

import "context"

// Command 对标 PSR Command 概念：一个可执行单元。
type Command interface {
	// Name 命令名，例如 "dev" / "modgen" / "bizgen" / "user:sync".
	Name() string
	// Description 一句话描述，用于 -h 输出。
	Description() string
	// Run 执行命令。args 是子命令后的参数（不包含命令名本身）。
	Run(ctx context.Context, args []string) error
}

// CommandProvider 由业务模块实现，提供一组命令。
// 例如 internal/user 实现它来暴露 user:* 系列命令。
type CommandProvider interface {
	Commands() []Command
}
