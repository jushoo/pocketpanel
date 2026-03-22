package models

import (
	"time"

	"gorm.io/gorm"
)

type Version struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	ServerType ServerType     `json:"server_type" gorm:"uniqueIndex:idx_version_type;not null"`
	Version    string         `json:"version" gorm:"uniqueIndex:idx_version_type;not null"`
	IsLatest   bool           `json:"is_latest" gorm:"default:false"`
	SyncedAt   time.Time      `json:"synced_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}
