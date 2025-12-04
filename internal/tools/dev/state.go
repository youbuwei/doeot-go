package dev

import (
	"os/exec"
	"sync"
	"time"
)

var (
	mu             sync.Mutex
	processes      = map[string]*exec.Cmd{}
	needGenOnce    bool
	lastRestart    time.Time
	lastGenerate   time.Time
	lastChange     string
	servicesGlobal []string
	devHTTPAddr    string
)
