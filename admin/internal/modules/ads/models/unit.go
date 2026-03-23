package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Unit struct {
	gorm.Model

	Title string
	Price int

	NetworkID int
	Network   Network

	PlacementID int
	Placement   Placement

	Data   string
	Active bool

	ArchivedAt *time.Time
}

func (original Unit) Archive(db *gorm.DB, archive *time.Time) error {
	original.ArchivedAt = archive

	if err := db.Save(&original).Error; err != nil {
		return fmt.Errorf("failed to archive: %w", err)
	}

	return nil
}
