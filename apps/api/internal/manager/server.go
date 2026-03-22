package manager

import (
	"fmt"

	"pocketpanel/api/internal/models"
)

// ServerManager orchestrates all server-related operations
type ServerManager struct {
	jarMgr     *JARManager
	processMgr *ProcessManager
	console    *FileConsoleProvider
}

// NewServerManager creates a new ServerManager
func NewServerManager() *ServerManager {
	return &ServerManager{
		jarMgr:     NewJARManager(BasePath),
		processMgr: NewProcessManager(),
		console:    NewFileConsoleProvider(),
	}
}

// PrepareServer prepares a server for starting: creates directories, downloads JAR, generates config
func (sm *ServerManager) PrepareServer(server *models.Server) error {
	// Ensure server directory exists
	serverDir, err := sm.jarMgr.EnsureServerDir(server.ID)
	if err != nil {
		return fmt.Errorf("failed to create server directory: %w", err)
	}

	// Download JAR if missing (only for vanilla servers currently)
	if server.Type == models.ServerTypeVanilla {
		if err := sm.jarMgr.DownloadIfMissing(server.ID, server.Version); err != nil {
			return fmt.Errorf("failed to download JAR: %w", err)
		}
	}

	// Generate server.properties
	if err := GenerateServerProperties(server, serverDir); err != nil {
		return fmt.Errorf("failed to generate server.properties: %w", err)
	}

	// Accept EULA
	if err := AcceptEULA(serverDir); err != nil {
		return fmt.Errorf("failed to accept EULA: %w", err)
	}

	return nil
}

// StartServer starts a Minecraft server
func (sm *ServerManager) StartServer(server *models.Server) error {
	// Check if already running
	if sm.processMgr.IsRunning(server.ID) {
		return fmt.Errorf("server %d is already running", server.ID)
	}

	// Prepare server (download JAR, generate config)
	if err := sm.PrepareServer(server); err != nil {
		return err
	}

	// Get JAR path
	jarPath := sm.jarMgr.GetServerJARPath(server.ID)

	// Start the process
	if err := sm.processMgr.Start(server, jarPath); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// StopServer stops a Minecraft server
func (sm *ServerManager) StopServer(serverID uint, force bool) error {
	return sm.processMgr.Stop(serverID, force)
}

// GetServerStatus returns the status of a server
func (sm *ServerManager) GetServerStatus(serverID uint) (*ServerStatus, error) {
	pid, running := sm.processMgr.GetPID(serverID)
	
	status := &ServerStatus{
		Running: running,
	}
	
	if running {
		status.PID = pid
	}

	return status, nil
}

// ServerStatus represents the status of a server
type ServerStatus struct {
	Running bool `json:"running"`
	PID     int  `json:"pid,omitempty"`
	Port    uint `json:"port,omitempty"`
}

// SubscribeConsole subscribes to console output for a server
func (sm *ServerManager) SubscribeConsole(serverID uint) (<-chan string, error) {
	return sm.console.Subscribe(serverID)
}

// UnsubscribeConsole unsubscribes from console output for a server
func (sm *ServerManager) UnsubscribeConsole(serverID uint) {
	sm.console.Unsubscribe(serverID)
}

// GetConsoleHistory returns recent console output for a server
func (sm *ServerManager) GetConsoleHistory(serverID uint, lines int) ([]string, error) {
	return sm.console.GetHistory(serverID, lines)
}

// CleanupDeadProcesses removes stale process entries
func (sm *ServerManager) CleanupDeadProcesses() {
	sm.processMgr.CleanupDeadProcesses()
}

// IsRunning checks if a server is currently running
func (sm *ServerManager) IsRunning(serverID uint) bool {
	return sm.processMgr.IsRunning(serverID)
}
