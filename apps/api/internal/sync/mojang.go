package sync

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"pocketpanel/api/internal/models"
)

const mojangManifestURL = "https://launchermeta.mojang.com/mc/game/version_manifest.json"

// versionManifestResponse represents the full version manifest from Mojang
type versionManifestResponse struct {
	Versions []mojangVersion `json:"versions"`
}

// mojangVersion represents a version entry in the manifest
type mojangVersion struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

// versionDetailResponse represents the detailed version info from Mojang
type versionDetailResponse struct {
	ID          string `json:"id"`
	Downloads   downloads `json:"downloads"`
	JavaVersion javaVersion `json:"java-version"`
}

type downloads struct {
	Server serverDownload `json:"server"`
}

type serverDownload struct {
	SHA1 string `json:"sha1"`
	URL  string `json:"url"`
}

type javaVersion struct {
	MajorVersion int `json:"majorVersion"`
}

// MojangFetcher handles fetching version info and downloading server JARs from Mojang.
type MojangFetcher struct {
	httpClient *http.Client
	// Cache version manifest to avoid repeated requests
	manifestCache *versionManifestResponse
	manifestTime  time.Time
}

func NewMojangFetcher() *MojangFetcher {
	return &MojangFetcher{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (m *MojangFetcher) ServerType() models.ServerType {
	return models.ServerTypeVanilla
}

// GetVersionManifest fetches the version manifest from Mojang with caching
func (m *MojangFetcher) GetVersionManifest() (*versionManifestResponse, error) {
	// Use cache if less than 1 hour old
	if m.manifestCache != nil && time.Since(m.manifestTime) < time.Hour {
		return m.manifestCache, nil
	}

	resp, err := m.httpClient.Get(mojangManifestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch manifest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var manifest versionManifestResponse
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return nil, fmt.Errorf("failed to decode manifest: %w", err)
	}

	m.manifestCache = &manifest
	m.manifestTime = time.Now()
	return &manifest, nil
}

// GetDownloadURL returns the download URL for a specific Minecraft version
func (m *MojangFetcher) GetDownloadURL(version string) (string, error) {
	manifest, err := m.GetVersionManifest()
	if err != nil {
		return "", err
	}

	// Find the version URL
	var versionURL string
	for _, v := range manifest.Versions {
		if v.ID == version {
			versionURL = v.URL
			break
		}
	}

	if versionURL == "" {
		return "", fmt.Errorf("version %s not found in manifest", version)
	}

	// Fetch the version detail
	resp, err := m.httpClient.Get(versionURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch version detail: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status for version detail: %d", resp.StatusCode)
	}

	var detail versionDetailResponse
	if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		return "", fmt.Errorf("failed to decode version detail: %w", err)
	}

	return detail.Downloads.Server.URL, nil
}

// DownloadJAR downloads the server JAR for the specified version to destPath
func (m *MojangFetcher) DownloadJAR(version string, destPath string) error {
	url, err := m.GetDownloadURL(version)
	if err != nil {
		return fmt.Errorf("failed to get download URL: %w", err)
	}

	resp, err := m.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download JAR: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status downloading JAR: %d", resp.StatusCode)
	}

	// Create parent directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create the file
	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy the content
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write JAR: %w", err)
	}

	return nil
}

// FetchVersions returns a list of versions from the Mojang manifest
func (m *MojangFetcher) FetchVersions() ([]models.Version, error) {
	manifest, err := m.GetVersionManifest()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	versions := make([]models.Version, 0)
	var latestRelease string

	for _, v := range manifest.Versions {
		if v.Type == "release" && IsValidSemver(v.ID) {
			if latestRelease == "" {
				latestRelease = v.ID
			}
			versions = append(versions, models.Version{
				ServerType: models.ServerTypeVanilla,
				Version:    v.ID,
				IsLatest:   v.ID == latestRelease,
				SyncedAt:   now,
			})
		}
	}

	return versions, nil
}
