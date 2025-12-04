package dev

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// 递归添加目录监听。
func addWatchRecursive(w *fsnotify.Watcher, root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}

		base := filepath.Base(path)
		// 跳过 vendor/.git/node_modules 等目录。
		if strings.HasPrefix(base, ".") || base == "vendor" || base == "node_modules" {
			return filepath.SkipDir
		}

		if err := w.Add(path); err != nil {
			return fmt.Errorf("watch %s: %w", path, err)
		}
		return nil
	})
}

// 是否需要触发重启（只关心 .go 文件）。
func shouldTrigger(ev fsnotify.Event) bool {
	if !strings.HasSuffix(ev.Name, ".go") {
		return false
	}

	// 忽略 bizgen 刚刚生成的 HTTP/RPC wrapper（zz_*.go），避免 go generate 后再多跑一轮。
	if isGeneratedWrapper(ev.Name) && time.Since(lastGenerate) < 2*time.Second {
		return false
	}

	if ev.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename|fsnotify.Remove) != 0 {
		return true
	}
	return false
}

// 是否需要 go generate：只在注解相关文件改动时置 true。
func needsGenerate(path string) bool {
	// 约定：注解 endpoint 都放在 interfaces/endpoint 下，或者以 _endpoint.go 结尾。
	if strings.Contains(path, "interfaces/endpoint") {
		return true
	}
	if strings.HasSuffix(path, "_endpoint.go") {
		return true
	}
	return false
}

// 判断是否是 bizgen 生成的 HTTP/RPC wrapper。
func isGeneratedWrapper(path string) bool {
	if strings.Contains(path, "/interfaces/http/") && strings.Contains(path, "zz_") {
		return true
	}
	if strings.Contains(path, "/interfaces/rpc/") && strings.Contains(path, "zz_") {
		return true
	}
	return false
}

// 处理 watcher 事件，将“重启请求”写入 restartCh。
func handleWatcherEvents(watcher *fsnotify.Watcher, restartCh chan<- struct{}) {
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

			// 非阻塞写入，避免事件风暴。
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
}

// 防抖 300ms，合并一波变化后再重启。
func debounceRestart(restartCh <-chan struct{}, services []string) {
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
}
