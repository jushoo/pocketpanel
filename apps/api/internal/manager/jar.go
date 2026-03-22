package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"pocketpanel/api/internal/sync"
)

const (
	// BasePath is the root directory for all server files
	BasePath = "/opt/pocketpanel/servers"
	// JARName is the name of the server JAR file
	JARName = "server.jar"
)

// JARManager handles downloading and caching server JARs.
type JARManager struct {
	basePath   string
	fetcher    *sync.MojangFetcher
}

// NewJARManager creates a new JARManager.
func NewJARManager(basePath string) *JARManager {
	if basePath == "" {
		basePath = BasePath
	}
	return &JARManager{
		basePath: basePath,
		fetcher:  sync.NewMojangFetcher(),
	}
}

// EnsureServerDir creates the server directory if it doesn't exist.
func (j *JARManager) EnsureServerDir(serverID uint) (string, error) {
	serverDir := j.GetServerDir(serverID)
	if err := os.MkdirAll(serverDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create server directory: %w", err)
	}
	return serverDir, nil
}

// GetServerDir returns the path to a server's directory.
func (j *JARManager) GetServerDir(serverID uint) string {
	return filepath.Join(j.basePath, fmt.Sprintf("%d", serverID))
}

// GetServerJARPath returns the full path to a server's JAR file.
func (j *JARManager) GetServerJARPath(serverID uint) string {
	return filepath.Join(j.GetServerDir(serverID), JARName)
}

// JARExists checks if the server JAR exists.
func (j *JARManager) JARExists(serverID uint) bool {
	jarPath := j.GetServerJARPath(serverID)
	_, err := os.Stat(jarPath)
	return err == nil
}

// DownloadIfMissing downloads the server JAR if it doesn't exist.
func (j *JARManager) DownloadIfMissing(serverID uint, version string) error {
	jarPath := j.GetServerJARPath(serverID)

	// Skip if JAR already exists
	if j.JARExists(serverID) {
		return nil
	}

	// Ensure directory exists
	_, err := j.EnsureServerDir(serverID)
	if err != nil {
		return err
	}

	// Download the JAR
	if err := j.fetcher.DownloadJAR(version, jarPath); err != nil {
		return fmt.Errorf("failed to download JAR for version %s: %w", version, err)
	}

	return nil
}

// Download downloads the server JAR, overwriting any existing JAR.
func (j *JARManager) Download(serverID uint, version string) error {
	// Ensure directory exists
	_, err := j.EnsureServerDir(serverID)
	if err != nil {
		return err
	}

	jarPath := j.GetServerJARPath(serverID)

	// Download the JAR
	if err := j.fetcher.DownloadJAR(version, jarPath); err != nil {
		return fmt.Errorf("failed to download JAR for version %s: %w", version, err)
	}

	return nil
}
