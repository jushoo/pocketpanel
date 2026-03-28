package fabric

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"pocketpanel/api/internal/models"
	"pocketpanel/api/internal/sync"
)

const fabricMetaURL = "https://meta.fabricmc.net/v2/versions/loader"

type FabricDownloader struct {
	httpClient *http.Client
}

func NewFabricDownloader() *FabricDownloader {
	return &FabricDownloader{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (f *FabricDownloader) ServerType() models.ServerType {
	return models.ServerTypeFabric
}

func (f *FabricDownloader) DownloadJAR(version string, destPath string) error {
	url := fmt.Sprintf("%s/%s/server/jar", fabricMetaURL, version)

	resp, err := f.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download Fabric JAR: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status downloading Fabric JAR: %d", resp.StatusCode)
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

type FabricVersion struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
}

func FetchVersions() ([]models.Version, error) {
	resp, err := http.Get("https://meta.fabricmc.net/v2/versions/game")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var fabricVersions []FabricVersion
	if err := json.NewDecoder(resp.Body).Decode(&fabricVersions); err != nil {
		return nil, err
	}

	now := time.Now()
	versions := make([]models.Version, 0)
	latestStable := ""

	for _, v := range fabricVersions {
		if v.Stable && sync.IsValidSemver(v.Version) {
			if latestStable == "" {
				latestStable = v.Version
			}
			versions = append(versions, models.Version{
				ServerType: models.ServerTypeFabric,
				Version:    v.Version,
				IsLatest:   v.Stable && v.Version == latestStable,
				SyncedAt:   now,
			})
		}
	}

	return versions, nil
}
