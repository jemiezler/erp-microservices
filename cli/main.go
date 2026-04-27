package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"slices"
	"sync"
	"syscall"
)

const (
	ColorReset  = "\033[0m"
	ColorBlue   = "\033[1;34m"
	ColorCyan   = "\033[0;36m"
	ColorYellow = "\033[1;33m"
	ColorGreen  = "\033[0;32m"
	ColorRed    = "\033[1;31m"
)

func logSystem(level, message string) {
	color := ColorBlue
	switch level {
	case "SUCCESS":
		color = ColorGreen
	case "ERROR":
		color = ColorRed
	}
	fmt.Printf("%s[SYSTEM] %s: %s%s\n", color, level, message, ColorReset)
}

func logDeps(action, target string) {
	fmt.Printf("%s[DEPS] %s: %s%s\n", ColorCyan, action, target, ColorReset)
}

func logLaunch(service, port, status string) {
	fmt.Printf("%s[LAUNCH] SERVICE: %-15s | PORT: %-5s | STATUS: %s%s\n", ColorYellow, service, port, status, ColorReset)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cli [run|dev|add|install]")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "run", "dev":
		runWorkspace(os.Args[2:], command == "dev")
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: cli add <service-name>")
			os.Exit(1)
		}
		addService(os.Args[2])
	case "install":
		installDeps(os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func installDeps(_ []string) {
	logSystem("INFO", "Validating workspace dependencies...")
	var targets []string
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err == nil && info.Name() == "go.mod" {
			targets = append(targets, filepath.Dir(path))
		}
		return nil
	})

	var wg sync.WaitGroup
	for _, target := range targets {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			logDeps("TIDY", t)
			cmd := exec.Command("go", "mod", "tidy")
			cmd.Dir = t
			cmd.Run()
		}(target)
	}
	wg.Wait()
	logSystem("SUCCESS", "Workspace dependencies synchronized.")
}

func runWorkspace(requestedServices []string, isDev bool) {
	if isDev {
		logSystem("INFO", "Initializing Development Mode...")
		installDeps(requestedServices)
	}

	logSystem("INFO", "Starting ERP Orchestrator...")
	// Try docker-compose then docker compose
	if err := exec.Command("docker-compose", "up", "-d").Run(); err != nil {
		exec.Command("docker", "compose", "up", "-d").Run()
	}

	var wg sync.WaitGroup
	ctxSignals := make(chan os.Signal, 1)
	signal.Notify(ctxSignals, os.Interrupt, syscall.SIGTERM)

	processes := []*os.Process{}
	var procMu sync.Mutex

	startProcess := func(dir, name, port string, command string, args ...string) {
		wg.Go(func() {

			// Resolve command for Windows compatibility
			fullCmd, err := exec.LookPath(command)
			if err != nil {
				fullCmd = command
			}

			cmd := exec.Command(fullCmd, args...)
			cmd.Dir = dir
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			logLaunch(name, port, "STARTING")
			if err := cmd.Start(); err != nil {
				logSystem("ERROR", fmt.Sprintf("Failed to start %s: %v", name, err))
				return
			}

			procMu.Lock()
			processes = append(processes, cmd.Process)
			procMu.Unlock()

			cmd.Wait()
		})
	}

	// 1. Start Backend (Using air or go run)
	startProcess("apis/api-gateway", "API Gateway", "8080", "go", "run", "main.go")

	servicesDir := filepath.Join("apis", "services")
	ports := map[string]string{
		"hr-service":      "8081",
		"finance-service": "8082",
	}

	entries, _ := os.ReadDir(servicesDir)
	runAll := len(requestedServices) == 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if runAll || contains(requestedServices, name) {
			// Corrected path to cmd/api
			servicePath := filepath.Join(servicesDir, name)
			startProcess(servicePath, name, ports[name], "go", "run", "./cmd/api")
		}
	}

	// 2. Start Frontend with CLEAN TURBO LOGS
	if runAll || contains(requestedServices, "frontend") {
		startProcess(
			filepath.Join("frontend", "apps", "host-app"),
			"host-app",
			"3000",
			"pnpm",
			"dev",
		)

		startProcess(
			filepath.Join("frontend", "apps", "hr-mfe"),
			"hr-mfe",
			"3001",
			"pnpm",
			"dev",
		)
	}

	logSystem("SUCCESS", "All services operational. Press CTRL+C to terminate.")

	<-ctxSignals
	fmt.Println()
	logSystem("INFO", "Shutting down all processes...")

	procMu.Lock()
	for _, proc := range processes {
		if proc != nil {
			if runtime.GOOS == "windows" {
				exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprint(proc.Pid)).Run()
			} else {
				proc.Signal(os.Interrupt)
			}
		}
	}
	procMu.Unlock()
	os.Exit(0)
}

func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}

func addService(name string) {
	logSystem("INFO", fmt.Sprintf("Scaffolding new service: %s", name))
}
