package sync

import (
	"log"

	"gorm.io/gorm"

	"pocketpanel/api/internal/models"
)

type UpstreamFetcher interface {
	FetchVersions() ([]models.Version, error)
	ServerType() models.ServerType
}

type Syncer struct {
	db       *gorm.DB
	fetchers map[models.ServerType]UpstreamFetcher
}

func NewSyncer(db *gorm.DB, fetchers ...UpstreamFetcher) *Syncer {
	fetcherMap := make(map[models.ServerType]UpstreamFetcher)
	for _, f := range fetchers {
		fetcherMap[f.ServerType()] = f
	}
	return &Syncer{
		db:       db,
		fetchers: fetcherMap,
	}
}

func (s *Syncer) SyncAll() error {
	for serverType, fetcher := range s.fetchers {
		if err := s.syncServerType(serverType, fetcher); err != nil {
			log.Printf("Failed to sync %s versions: %v", serverType, err)
		}
	}
	return nil
}

func (s *Syncer) SyncServerType(serverType models.ServerType) error {
	fetcher, ok := s.fetchers[serverType]
	if !ok {
		return nil
	}
	return s.syncServerType(serverType, fetcher)
}

func (s *Syncer) syncServerType(serverType models.ServerType, fetcher UpstreamFetcher) error {
	versions, err := fetcher.FetchVersions()
	if err != nil {
		return err
	}

	if len(versions) == 0 {
		return nil
	}

	seen := make(map[string]bool)
	deduped := make([]models.Version, 0, len(versions))
	for _, v := range versions {
		if !seen[v.Version] {
			seen[v.Version] = true
			deduped = append(deduped, v)
		}
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("server_type = ?", serverType).Delete(&models.Version{}).Error; err != nil {
			return err
		}
		return tx.Create(&deduped).Error
	})
}
