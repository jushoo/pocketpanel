package manager

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pocketpanel/api/internal/sync/vanilla"

	stdlib_sync "sync"
)

// ConsoleProvider defines the interface for console output streaming
type ConsoleProvider interface {
	// Subscribe returns a channel that receives console output lines
	Subscribe(serverID uint) (<-chan string, error)
	// Unsubscribe stops the subscription
	Unsubscribe(serverID uint)
	// GetHistory returns recent console output
	GetHistory(serverID uint, lines int) ([]string, error)
}

// FileConsoleProvider implements ConsoleProvider using file-based logging
type FileConsoleProvider struct {
	mu       stdlib_sync.RWMutex
	subs     map[uint][]chan string
	stopChan map[uint]chan struct{}
}

// NewFileConsoleProvider creates a new FileConsoleProvider
func NewFileConsoleProvider() *FileConsoleProvider {
	return &FileConsoleProvider{
		subs:     make(map[uint][]chan string),
		stopChan: make(map[uint]chan struct{}),
	}
}

// Subscribe returns a channel that receives console output lines from the log file
func (fcp *FileConsoleProvider) Subscribe(serverID uint) (<-chan string, error) {
	fcp.mu.Lock()
	defer fcp.mu.Unlock()

	logPath := fcp.getLogPath(serverID)

	// Create channel
	ch := make(chan string, 100) // Buffer to prevent blocking

	// Add to subscriptions
	fcp.subs[serverID] = append(fcp.subs[serverID], ch)

	// Create stop channel for this subscription
	if fcp.stopChan[serverID] == nil {
		fcp.stopChan[serverID] = make(chan struct{})
	}

	// Start tail goroutine
	go fcp.tailLog(serverID, ch, logPath)

	return ch, nil
}

// Unsubscribe stops all subscriptions for a server
func (fcp *FileConsoleProvider) Unsubscribe(serverID uint) {
	fcp.mu.Lock()
	defer fcp.mu.Unlock()

	// Close stop channel
	if stopCh, ok := fcp.stopChan[serverID]; ok {
		close(stopCh)
		delete(fcp.stopChan, serverID)
	}

	// Close all subscriber channels
	if subs, ok := fcp.subs[serverID]; ok {
		for _, ch := range subs {
			close(ch)
		}
		delete(fcp.subs, serverID)
	}
}

// GetHistory returns recent console output from the log file
func (fcp *FileConsoleProvider) GetHistory(serverID uint, lines int) ([]string, error) {
	logPath := fcp.getLogPath(serverID)

	file, err := os.Open(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	// Read all lines
	var allLines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		allLines = append(allLines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read log file: %w", err)
	}

	// Return last n lines
	if lines > len(allLines) {
		lines = len(allLines)
	}
	if lines == 0 {
		return []string{}, nil
	}

	start := len(allLines) - lines
	return allLines[start:], nil
}

func (fcp *FileConsoleProvider) getLogPath(serverID uint) string {
	jarMgr := NewJARManager("", vanilla.NewMojangDownloader())
	serverDir := jarMgr.GetServerDir(serverID)
	return filepath.Join(serverDir, ConsoleLogName)
}

func (fcp *FileConsoleProvider) tailLog(serverID uint, ch chan string, logPath string) {
	file, err := os.Open(logPath)
	if err != nil {
		if !os.IsNotExist(err) {
			ch <- fmt.Sprintf("[ERROR] Failed to open log file: %v", err)
		}
		close(ch)
		return
	}
	defer file.Close()

	// Seek to end of file
	file.Seek(0, io.SeekEnd)

	reader := bufio.NewReader(file)

	// Get stop channel
	fcp.mu.RLock()
	stopCh := fcp.stopChan[serverID]
	fcp.mu.RUnlock()

	for {
		select {
		case <-stopCh:
			close(ch)
			return
		default:
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				// No more data, wait a bit and retry
				time.Sleep(100 * time.Millisecond)
				continue
			}
			if err != nil {
				ch <- fmt.Sprintf("[ERROR] Failed to read log: %v", err)
				close(ch)
				return
			}
			// Remove trailing newline
			line = strings.TrimSuffix(line, "\n")
			if line != "" {
				select {
				case ch <- line:
				case <-stopCh:
					close(ch)
					return
				}
			}
		}
	}
}

// ServerConsole wraps the server management with console capabilities
type ServerConsole struct {
	processMgr *ProcessManager
	jarMgr     *JARManager
	console    ConsoleProvider
}

// NewServerConsole creates a new ServerConsole
func NewServerConsole() *ServerConsole {
	return &ServerConsole{
		processMgr: NewProcessManager(),
		jarMgr:     NewJARManager("", vanilla.NewMojangDownloader()),
		console:    NewFileConsoleProvider(),
	}
}

// GetProcessManager returns the process manager
func (sc *ServerConsole) GetProcessManager() *ProcessManager {
	return sc.processMgr
}

// GetJARManager returns the JAR manager
func (sc *ServerConsole) GetJARManager() *JARManager {
	return sc.jarMgr
}

// SubscribeConsole subscribes to console output for a server
func (sc *ServerConsole) SubscribeConsole(serverID uint) (<-chan string, error) {
	return sc.console.Subscribe(serverID)
}

// UnsubscribeConsole unsubscribes from console output
func (sc *ServerConsole) UnsubscribeConsole(serverID uint) {
	sc.console.Unsubscribe(serverID)
}

// GetConsoleHistory returns recent console output
func (sc *ServerConsole) GetConsoleHistory(serverID uint, lines int) ([]string, error) {
	return sc.console.GetHistory(serverID, lines)
}
