package sync

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"pocketpanel/api/internal/models"
)

const fabricMetaURL = "https://meta.fabricmc.net/v2/versions/game"

type FabricFetcher struct {
	httpClient *http.Client
}

func NewFabricFetcher() *FabricFetcher {
	return &FabricFetcher{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (f *FabricFetcher) ServerType() models.ServerType {
	return models.ServerTypeFabric
}

type fabricVersion struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
}

func (f *FabricFetcher) FetchVersions() ([]models.Version, error) {
	resp, err := f.httpClient.Get(fabricMetaURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var fabricVersions []fabricVersion
	if err := json.NewDecoder(resp.Body).Decode(&fabricVersions); err != nil {
		return nil, err
	}

	now := time.Now()
	versions := make([]models.Version, 0)
	latestStable := ""

	for _, v := range fabricVersions {
		if v.Stable && IsValidSemver(v.Version) {
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
