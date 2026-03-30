package sync

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

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
			return fmt.Errorf("sync failed for %s: %w", serverType, err)
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

	// Deduplicate by version string (normalize the data)
	// The unique constraint is on (server_type, version), so we need unique version strings
	seen := make(map[string]bool)
	deduped := make([]models.Version, 0, len(versions))
	for _, v := range versions {
		// Ensure ServerType is set correctly
		v.ServerType = serverType

		// Skip if we've already seen this version
		if seen[v.Version] {
			continue
		}
		seen[v.Version] = true
		deduped = append(deduped, v)
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Delete existing versions for this server type
		if err := tx.Where("server_type = ?", serverType).Delete(&models.Version{}).Error; err != nil {
			return err
		}

		// Insert with conflict handling (skip duplicates if any slip through)
		return tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "server_type"}, {Name: "version"}},
			DoNothing: true,
		}).Create(&deduped).Error
	})
}
