package models

import (
	"time"

	"gorm.io/gorm"
)

type ServerType string

const (
	ServerTypeVanilla ServerType = "vanilla"
	ServerTypeFabric  ServerType = "fabric"
)

type Server struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"uniqueIndex;not null"`
	Type      ServerType     `json:"type" gorm:"not null"`
	Version   string         `json:"version" gorm:"not null"`
	MinMem    uint           `json:"min_mem" gorm:"not null"`
	MaxMem    uint           `json:"max_mem" gorm:"not null"`
	Port      uint           `json:"port" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
