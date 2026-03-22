package sync

import (
	"encoding/json"
	"net/http"
	"time"

	"pocketpanel/api/internal/models"
)

const mojangManifestURL = "https://launchermeta.mojang.com/mc/game/version_manifest.json"

type MojangFetcher struct {
	httpClient *http.Client
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

type mojangManifest struct {
	Versions []mojangVersion `json:"versions"`
}

type mojangVersion struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

func (m *MojangFetcher) FetchVersions() ([]models.Version, error) {
	resp, err := m.httpClient.Get(mojangManifestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	var manifest mojangManifest
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
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
