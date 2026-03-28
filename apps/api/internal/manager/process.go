package manager

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"time"

	"pocketpanel/api/internal/models"
)

const (
	// DefaultShutdownTimeout is the default time to wait for graceful shutdown
	DefaultShutdownTimeout = 30 * time.Second
	// ConsoleLogName is the name of the console log file
	ConsoleLogName = "logs/latest.log"
)

// ProcessInfo contains information about a running server process
type ProcessInfo struct {
	PID     int
	Server  *models.Server
	StartAt time.Time
}

// ProcessManager manages Minecraft server processes.
type ProcessManager struct {
	mu              sync.RWMutex
	processes       map[uint]*ProcessInfo
	shutdownTimeout time.Duration
	serversPath     string
}

// NewProcessManager creates a new ProcessManager.
func NewProcessManager(serversPath string) *ProcessManager {
	return &ProcessManager{
		processes:       make(map[uint]*ProcessInfo),
		shutdownTimeout: DefaultShutdownTimeout,
		serversPath:     serversPath,
	}
}

// SetShutdownTimeout sets the timeout for graceful shutdown.
func (pm *ProcessManager) SetShutdownTimeout(timeout time.Duration) {
	pm.shutdownTimeout = timeout
}

// Start starts a Minecraft server process.
func (pm *ProcessManager) Start(server *models.Server, jarPath string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Check if already running
	if _, exists := pm.processes[server.ID]; exists {
		return fmt.Errorf("server %d is already running", server.ID)
	}

	// Check if JAR exists
	if _, err := os.Stat(jarPath); os.IsNotExist(err) {
		return fmt.Errorf("JAR file not found: %s", jarPath)
	}

	// Build JVM arguments
	jvmArgs := []string{
		"-Xms" + strconv.Itoa(int(server.MinMem)) + "M",
		"-Xmx" + strconv.Itoa(int(server.MaxMem)) + "M",
		"-jar",
		jarPath,
		"nogui",
	}

	// Create the command
	cmd := exec.Command("java", jvmArgs...)
	cmd.Dir = filepath.Dir(jarPath) // Set working directory to server folder

	// Create logs directory
	logsDir := filepath.Join(cmd.Dir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Set up log files
	stdoutFile, err := os.OpenFile(
		filepath.Join(logsDir, "stdout.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return fmt.Errorf("failed to open stdout log: %w", err)
	}

	stderrFile, err := os.OpenFile(
		filepath.Join(logsDir, "stderr.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		stdoutFile.Close()
		return fmt.Errorf("failed to open stderr log: %w", err)
	}

	cmd.Stdout = stdoutFile
	cmd.Stderr = stderrFile

	// Start the process
	if err := cmd.Start(); err != nil {
		stdoutFile.Close()
		stderrFile.Close()
		return fmt.Errorf("failed to start process: %w", err)
	}

	// Store process info
	pm.processes[server.ID] = &ProcessInfo{
		PID:     cmd.Process.Pid,
		Server:  server,
		StartAt: time.Now(),
	}

	// Clean up file handles (process inherits them)
	stdoutFile.Close()
	stderrFile.Close()

	return nil
}

// Stop stops a Minecraft server process.
func (pm *ProcessManager) Stop(serverID uint, force bool) error {
	pm.mu.Lock()
	info, exists := pm.processes[serverID]
	if !exists {
		return fmt.Errorf("server %d is not running", serverID)
	}

	proc, err := os.FindProcess(info.PID)
	pm.mu.Unlock()

	if err != nil {
		return fmt.Errorf("failed to find process: %w", err)
	}

	if force {
		// Send SIGKILL immediately
		if err := proc.Kill(); err != nil {
			return fmt.Errorf("failed to kill process: %w", err)
		}
	} else {
		// Send SIGINT for graceful shutdown
		if err := proc.Signal(syscall.SIGINT); err != nil {
			return fmt.Errorf("failed to send SIGINT: %w", err)
		}

		// Wait for graceful shutdown with timeout
		done := make(chan error, 1)
		go func() {
			_, err := proc.Wait()
			done <- err
		}()

		select {
		case <-time.After(pm.shutdownTimeout):
			// Force kill after timeout
			if err := proc.Kill(); err != nil {
				return fmt.Errorf("failed to force kill after timeout: %w", err)
			}
			return fmt.Errorf("server did not shut down gracefully, force killed")
		case err := <-done:
			if err != nil {
				return fmt.Errorf("process exited with error: %w", err)
			}
		}
	}

	// Clean up process map
	pm.mu.Lock()
	delete(pm.processes, serverID)
	pm.mu.Unlock()

	return nil
}

// GetPID returns the PID of a running server, or false if not running.
func (pm *ProcessManager) GetPID(serverID uint) (int, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	info, exists := pm.processes[serverID]
	if !exists {
		return 0, false
	}

	// Verify process is still running
	proc, err := os.FindProcess(info.PID)
	if err != nil || proc.Signal(syscall.Signal(0)) != nil {
		return 0, false
	}

	return info.PID, true
}

// IsRunning checks if a server is currently running.
func (pm *ProcessManager) IsRunning(serverID uint) bool {
	_, running := pm.GetPID(serverID)
	return running
}

// GetProcessInfo returns information about a running process.
func (pm *ProcessManager) GetProcessInfo(serverID uint) (*ProcessInfo, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	info, exists := pm.processes[serverID]
	if !exists {
		return nil, false
	}

	// Verify process is still running
	proc, err := os.FindProcess(info.PID)
	if err != nil || proc.Signal(syscall.Signal(0)) != nil {
		return nil, false
	}

	return info, true
}

// ListRunning returns a list of all running server IDs.
func (pm *ProcessManager) ListRunning() []uint {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	running := make([]uint, 0, len(pm.processes))
	for id := range pm.processes {
		running = append(running, id)
	}
	return running
}

// CleanupDeadProcesses removes stale process entries.
func (pm *ProcessManager) CleanupDeadProcesses() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for id, info := range pm.processes {
		proc, err := os.FindProcess(info.PID)
		if err == nil && proc.Signal(syscall.Signal(0)) == nil {
			continue // Process still running
		}
		delete(pm.processes, id)
	}
}

// WaitForShutdown blocks until the server process exits.
func (pm *ProcessManager) WaitForShutdown(serverID uint) error {
	pm.mu.RLock()
	info, exists := pm.processes[serverID]
	if !exists {
		pm.mu.RUnlock()
		return fmt.Errorf("server %d is not running", serverID)
	}
	pid := info.PID
	pm.mu.RUnlock()

	proc, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %w", err)
	}

	_, err = proc.Wait()
	return err
}

// SendCommand sends a command to the server's stdin.
// Note: This requires the server process to have stdin connected.
func (pm *ProcessManager) SendCommand(serverID uint, command string) error {
	pm.mu.RLock()
	_, exists := pm.processes[serverID]
	pm.mu.RUnlock()

	if !exists {
		return fmt.Errorf("server %d is not running", serverID)
	}

	pipePath := filepath.Join(pm.getServerDir(serverID), "logs", "stdin.pipe")

	f, err := os.OpenFile(pipePath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open stdin pipe: %w", err)
	}
	defer f.Close()

	_, err = f.WriteString(command + "\n")
	return err
}

// getServerDir returns the server directory path for the given server ID.
func (pm *ProcessManager) getServerDir(serverID uint) string {
	return filepath.Join(pm.serversPath, fmt.Sprintf("%d", serverID))
}
