package legalentities

import (
	"time"

	"gorm.io/gorm"
)

type LegalEntity struct {
	UUID      string         `gorm:"primaryKey"`
	Name      string         `gorm:"not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
