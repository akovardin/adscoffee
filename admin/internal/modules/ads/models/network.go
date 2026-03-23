package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Network struct {
	gorm.Model

	Title  string
	Name   string
	Data   string
	Active bool

	ArchivedAt *time.Time
}

func (original Network) Archive(db *gorm.DB, archive *time.Time) error {
	original.ArchivedAt = archive

	if err := db.Save(&original).Error; err != nil {
		return fmt.Errorf("failed to archive: %w", err)
	}

	return nil
}
