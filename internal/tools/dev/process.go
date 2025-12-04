package dev

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"
)

// 重启所有服务：可选是否先执行 go generate。
func restartAllServices(services []string, runGenerate bool) {
	mu.Lock()
	defer mu.Unlock()

	// 先关掉旧的。
	for name, cmd := range processes {
		if cmd != nil && cmd.Process != nil {
			log.Printf("dev: stopping %s...", name)
			killProcessTree(cmd)
			_, _ = cmd.Process.Wait()
			processes[name] = nil
		}
	}

	// 需要的话跑一轮 go generate。
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

	// 再启动所有服务。
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

		// 类 Unix：放到一个新的进程组，方便整体 Kill（go run + 真正服务进程）。
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

// 停掉所有服务。
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

// 杀掉 go run 以及它 fork 出来的整个进程组（避免端口占用）。
func killProcessTree(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}

	// Windows 简单粗暴，直接 Kill。
	if runtime.GOOS == "windows" {
		_ = cmd.Process.Kill()
		return
	}

	// 类 Unix：按进程组杀。
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err != nil {
		// 拿不到进程组，就退化成只杀自己。
		_ = cmd.Process.Kill()
		return
	}

	// 负的 pgid 表示杀整个进程组。
	_ = syscall.Kill(-pgid, syscall.SIGKILL)
}
