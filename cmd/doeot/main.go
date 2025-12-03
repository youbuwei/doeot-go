package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	sub := os.Args[1]
	// 顶层帮助
	if sub == "-h" || sub == "--help" || sub == "help" {
		printUsage()
		return
	}

	args := os.Args[2:]

	switch sub {
	case "dev":
		runSub("dev", "./cmd/dev", args)
	case "modgen":
		runSub("modgen", "./cmd/modgen", args)
	case "bizgen":
		runSub("bizgen", "./cmd/bizgen", args)
	default:
		fmt.Fprintf(os.Stderr, "doeot: unknown command %q\n\n", sub)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`doeot - DOEOT 项目开发工具集合

用法:
  doeot <command> [arguments]

可用命令:
  dev       启动开发模式（热更新、多服务、HTTP 面板）
  modgen    生成业务模块骨架 (domain/app/repo/endpoint/module + bizgen)
  bizgen    根据 endpoint 注解生成 HTTP/RPC 包装代码

示例:
  # 开发模式，同时跑 HTTP + RPC 服务
  doeot dev -services user-api,user-rpc -dev-http :18080

  # 生成一个名为 order 的新模块
  doeot modgen -name order

  # 仅针对某个模块手动生成 HTTP/RPC wrapper
  doeot bizgen -module user

说明:
  - 每个子命令本身都支持 -h/--help 查看自己的参数说明。
  - 如果你是通过 go run 运行，可使用:
      go run ./cmd/doeot dev -services user-api,user-rpc
`)
}

// 通过 go run 调用对应的子命令
func runSub(name, pkgPath string, args []string) {
	// 支持将来用二进制直接跑时，也能把子命令的 -h 转发过去
	fullArgs := append([]string{"run", pkgPath}, args...)
	fmt.Printf("doeot: exec go %s\n", strings.Join(fullArgs, " "))

	cmd := exec.Command("go", fullArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		// 如果是正常的 exit code（例如参数错误），就按子进程的 code 退出
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "doeot: failed to run %s: %v\n", name, err)
		os.Exit(1)
	}
}
