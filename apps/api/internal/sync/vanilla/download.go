package vanilla

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

type versionManifestResponse struct {
	Versions []mojangVersion `json:"versions"`
}

type mojangVersion struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

type versionDetailResponse struct {
	ID          string      `json:"id"`
	Downloads   downloads   `json:"downloads"`
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

type MojangDownloader struct {
	httpClient    *http.Client
	manifestCache *versionManifestResponse
	manifestTime  time.Time
}

func NewMojangDownloader() *MojangDownloader {
	return &MojangDownloader{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (m *MojangDownloader) ServerType() models.ServerType {
	return models.ServerTypeVanilla
}

func (m *MojangDownloader) GetVersionManifest() (*versionManifestResponse, error) {
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

func (m *MojangDownloader) GetDownloadURL(version string) (string, error) {
	manifest, err := m.GetVersionManifest()
	if err != nil {
		return "", err
	}

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

func (m *MojangDownloader) DownloadJAR(version string, destPath string) error {
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

	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write JAR: %w", err)
	}

	return nil
}
