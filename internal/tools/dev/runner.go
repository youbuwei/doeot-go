package dev

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

// Run 是 dev 工具的主入口。
func Run(ctx context.Context, cfg Config) error {
	if len(cfg.Services) == 0 {
		return fmt.Errorf("no services specified")
	}
	servicesGlobal = cfg.Services
	devHTTPAddr = cfg.HTTPPanelAddr

	// 先跑一次 go generate + 启动服务（此时还没有 watcher）。
	restartAllServices(cfg.Services, true)

	// 启动 HTTP 面板。
	go startHTTPPanel(devHTTPAddr)

	// 初始化 watcher（在第一次 restart 之后）。
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("dev: new watcher error: %w", err)
	}
	defer watcher.Close()

	// 监听目录。
	watchDirs := []string{"internal", "pkg", "cmd"}
	for _, d := range watchDirs {
		if err := addWatchRecursive(watcher, d); err != nil {
			log.Printf("dev: watch dir %s error: %v", d, err)
		}
	}

	// 文件变更事件 → 合并后触发重启。
	restartCh := make(chan struct{}, 1)

	go handleWatcherEvents(watcher, restartCh)
	go debounceRestart(restartCh, cfg.Services)

	// 信号 + 控制台命令。
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go runCLI(cfg.Services, sigCh)

	<-sigCh

	log.Println("dev: shutting down...")
	stopAllServices()
	return nil
}

// 控制台命令循环。
func runCLI(services []string, sigCh chan os.Signal) {
	reader := bufio.NewReader(os.Stdin)

	log.Printf("dev: running services: %s", services)
	log.Printf("dev: HTTP panel: http://localhost%s/", devHTTPAddr)
	log.Println("dev: commands: [r] restart (go generate + restart), [s] status, [q] quit")

	for {
		fmt.Print("dev> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				sigCh <- os.Interrupt
				return
			}
			log.Printf("dev: read stdin error: %v", err)
			continue
		}
		line = strings.TrimSpace(line)
		switch line {
		case "":
			continue
		case "r", "restart":
			log.Println("dev: manual restart (with go generate)")
			restartAllServices(services, true)
		case "s", "status":
			printStatus()
		case "h", "help", "?":
			log.Println("dev: commands: [r] restart (go generate + restart), [s] status, [q] quit")
		case "q", "quit", "exit":
			sigCh <- os.Interrupt
			return
		default:
			log.Printf("dev: unknown command %q (type 'h' for help)", line)
		}
	}
}

func printStatus() {
	mu.Lock()
	defer mu.Unlock()

	if len(processes) == 0 {
		log.Println("dev: no services running")
		return
	}
	log.Println("dev: services status:")
	for name, cmd := range processes {
		status := "stopped"
		if cmd != nil && cmd.Process != nil && cmd.ProcessState == nil {
			status = "running"
		}
		log.Printf("  - %s: %s", name, status)
	}
}
