package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	mu             sync.Mutex
	processes      = map[string]*exec.Cmd{} // serviceName -> process
	needGenOnce    bool                     // 本轮变更是否需要 go generate
	lastRestart    time.Time                // 最近一次重启所有服务的时间
	lastGenerate   time.Time                // 最近一次 go generate 的时间
	lastChange     string                   // 最近一次触发的文件变更
	servicesGlobal []string                 // 当前托管的服务列表
	devHTTPAddr    string                   // HTTP 面板地址
)

func main() {
	servicesFlag := flag.String("services", "user-api", "逗号分隔的服务名，例如: user-api,user-rpc")
	devHTTP := flag.String("dev-http", ":18080", "dev 工具 HTTP 面板监听地址，例如 :18080")
	flag.Parse()

	services := splitAndTrim(*servicesFlag)
	if len(services) == 0 {
		log.Fatal("dev: no services specified, use -services=user-api or -services=user-rpc,user-api")
	}
	servicesGlobal = services
	devHTTPAddr = *devHTTP

	// 先跑一次 go generate + 启动服务（此时还没有 watcher，不会收到 wrapper 的写事件）
	restartAllServices(services, true)

	// 再启动 HTTP 面板
	go startHTTPPanel(devHTTPAddr)

	// 初始化 watcher（注意：在第一次 restart 之后）
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("dev: new watcher error: %v", err)
	}
	defer watcher.Close()

	// 监听目录
	watchDirs := []string{"internal", "pkg", "cmd"}
	for _, d := range watchDirs {
		if err := addWatchRecursive(watcher, d); err != nil {
			log.Printf("dev: watch dir %s error: %v", d, err)
		}
	}

	// 文件变更事件 → 合并后触发重启
	restartCh := make(chan struct{}, 1)

	go func() {
		for {
			select {
			case ev, ok := <-watcher.Events:
				if !ok {
					return
				}
				if !shouldTrigger(ev) {
					continue
				}

				log.Printf("dev: file changed: %s (%s)", ev.Name, ev.Op.String())

				mu.Lock()
				if needsGenerate(ev.Name) {
					needGenOnce = true
				}
				lastChange = fmt.Sprintf("%s (%s)", ev.Name, ev.Op.String())
				mu.Unlock()

				// 非阻塞写入，避免事件风暴
				select {
				case restartCh <- struct{}{}:
				default:
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("dev: watcher error: %v", err)
			}
		}
	}()

	// 防抖 300ms，合并一波变化后再重启
	go func() {
		var timer *time.Timer

		for range restartCh {
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(300*time.Millisecond, func() {
				mu.Lock()
				gen := needGenOnce
				needGenOnce = false
				mu.Unlock()

				restartAllServices(services, gen)
			})
		}
	}()

	// 信号 + 控制台命令
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go runCLI(services, sigCh)

	<-sigCh

	log.Println("dev: shutting down...")
	stopAllServices()
}

// --- 控制台命令 ---

func runCLI(services []string, sigCh chan os.Signal) {
	reader := bufio.NewReader(os.Stdin)

	log.Printf("dev: running services: %s", strings.Join(services, ", "))
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

// --- 重启/停止服务 ---

// 重启所有服务：可选是否先执行 go generate
func restartAllServices(services []string, runGenerate bool) {
	mu.Lock()
	defer mu.Unlock()

	// 先关掉旧的
	for name, cmd := range processes {
		if cmd != nil && cmd.Process != nil {
			log.Printf("dev: stopping %s...", name)
			killProcessTree(cmd)
			_, _ = cmd.Process.Wait()
			processes[name] = nil
		}
	}

	// 需要的话跑一轮 go generate
	if runGenerate {
		log.Println("dev: go generate ./...")
		gen := exec.Command("go", "generate", "./...")
		gen.Stdout = os.Stdout
		gen.Stderr = os.Stderr
		if err := gen.Run(); err != nil {
			log.Printf("dev: go generate error: %v", err)
		} else {
			lastGenerate = time.Now()
		}
	}

	// 再启动所有服务
	for _, svc := range services {
		svc := svc
		if svc == "" {
			continue
		}
		servicePath := fmt.Sprintf("./cmd/%s", svc)
		log.Printf("dev: go run %s", servicePath)

		cmd := exec.Command("go", "run", servicePath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// 类 Unix：放到一个新的进程组，方便整体 Kill（go run + 真正服务进程）
		if runtime.GOOS != "windows" {
			cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		}

		if err := cmd.Start(); err != nil {
			log.Printf("dev: start service %s error: %v", svc, err)
			continue
		}

		processes[svc] = cmd

		go func(name string, c *exec.Cmd) {
			if err := c.Wait(); err != nil {
				log.Printf("dev: service %s exited: %v", name, err)
			} else {
				log.Printf("dev: service %s exited normally", name)
			}
		}(svc, cmd)
	}

	lastRestart = time.Now()
}

// 停掉所有服务
func stopAllServices() {
	mu.Lock()
	defer mu.Unlock()

	for name, cmd := range processes {
		if cmd != nil && cmd.Process != nil {
			log.Printf("dev: stopping %s...", name)
			killProcessTree(cmd)
			_, _ = cmd.Process.Wait()
			processes[name] = nil
		}
	}
}

// 杀掉 go run 以及它 fork 出来的整个进程组（避免端口占用）
func killProcessTree(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}

	// Windows 简单粗暴，直接 Kill
	if runtime.GOOS == "windows" {
		_ = cmd.Process.Kill()
		return
	}

	// 类 Unix：按进程组杀
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err != nil {
		// 拿不到进程组，就退化成只杀自己
		_ = cmd.Process.Kill()
		return
	}

	// 负的 pgid 表示杀整个进程组
	_ = syscall.Kill(-pgid, syscall.SIGKILL)
}

// --- 文件监听相关 ---

// 递归添加目录监听
func addWatchRecursive(w *fsnotify.Watcher, root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}

		base := filepath.Base(path)
		// 跳过 vendor/.git/node_modules 等目录
		if strings.HasPrefix(base, ".") || base == "vendor" || base == "node_modules" {
			return filepath.SkipDir
		}

		if err := w.Add(path); err != nil {
			return fmt.Errorf("watch %s: %w", path, err)
		}
		return nil
	})
}

// 是否需要触发重启（只关心 .go 文件）
func shouldTrigger(ev fsnotify.Event) bool {
	if !strings.HasSuffix(ev.Name, ".go") {
		return false
	}

	// 忽略 bizgen 刚刚生成的 HTTP/RPC wrapper（zz_*.go），避免 go generate 后再多跑一轮
	if isGeneratedWrapper(ev.Name) && time.Since(lastGenerate) < 2*time.Second {
		return false
	}

	if ev.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename|fsnotify.Remove) != 0 {
		return true
	}
	return false
}

// 是否需要 go generate：只在注解相关文件改动时置 true
func needsGenerate(path string) bool {
	// 约定：注解 endpoint 都放在 interfaces/endpoint 下，或者以 _endpoint.go 结尾
	if strings.Contains(path, "interfaces/endpoint") {
		return true
	}
	if strings.HasSuffix(path, "_endpoint.go") {
		return true
	}
	return false
}

// 判断是否是 bizgen 生成的 HTTP/RPC wrapper
func isGeneratedWrapper(path string) bool {
	if strings.Contains(path, "/interfaces/http/") && strings.Contains(path, "zz_") {
		return true
	}
	if strings.Contains(path, "/interfaces/rpc/") && strings.Contains(path, "zz_") {
		return true
	}
	return false
}

// 工具函数：切分 services
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

// --- HTTP 状态面板 ---

func startHTTPPanel(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleStatus)

	log.Printf("dev: HTTP panel listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Printf("dev: HTTP panel error: %v", err)
	}
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintln(w, "<!doctype html><html><head><meta charset='utf-8'><title>dev panel</title>")
	fmt.Fprintln(w, "<style>body{font-family:-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica,Arial,sans-serif;margin:16px;} table{border-collapse:collapse;margin-top:8px;} th,td{border:1px solid #ddd;padding:4px 8px;} th{background:#f5f5f5;}</style></head><body>")
	fmt.Fprintln(w, "<h1>Dev Panel</h1>")
	fmt.Fprintf(w, "<p><b>Services:</b> %s</p>", strings.Join(servicesGlobal, ", "))
	fmt.Fprintf(w, "<p><b>Last restart:</b> %s</p>", formatTime(lastRestart))
	fmt.Fprintf(w, "<p><b>Last go generate:</b> %s</p>", formatTime(lastGenerate))
	fmt.Fprintf(w, "<p><b>Last file change:</b> %s</p>", lastChange)

	fmt.Fprintln(w, "<h2>Processes</h2>")
	fmt.Fprintln(w, "<table><tr><th>Name</th><th>Status</th><th>PID</th></tr>")
	for name, cmd := range processes {
		status := "stopped"
		pid := "-"
		if cmd != nil && cmd.Process != nil && cmd.ProcessState == nil {
			status = "running"
			pid = fmt.Sprintf("%d", cmd.Process.Pid)
		}
		fmt.Fprintf(w, "<tr><td>%s</td><td>%s</td><td>%s</td></tr>", name, status, pid)
	}
	fmt.Fprintln(w, "</table>")
	fmt.Fprintln(w, "<p>Commands (in terminal): r=restart (go generate + restart), s=status, q=quit</p>")
	fmt.Fprintln(w, "</body></html>")
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format("2006-01-02 15:04:05")
}
