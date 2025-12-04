package dev

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// 启动 HTTP 状态面板。
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
