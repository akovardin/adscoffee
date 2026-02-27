package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Placement struct {
	gorm.Model

	Name       string
	ArchivedAt *time.Time
	Active     bool
}

func (original Placement) Archive(db *gorm.DB, archive *time.Time) error {
	original.ArchivedAt = archive

	if err := db.Save(&original).Error; err != nil {
		return fmt.Errorf("failed to archive: %w", err)
	}

	return nil
}
