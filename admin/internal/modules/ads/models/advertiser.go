package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Advertiser struct {
	gorm.Model

	Title  string
	Info   string
	Active bool

	Start time.Time
	End   time.Time

	Targeting string
	Budget    string
	Capping   string
	Timetable string

	OrdContract string

	ArchivedAt *time.Time
}

func (original Advertiser) Archive(db *gorm.DB, archive *time.Time) error {
	original.ArchivedAt = archive

	if err := db.Save(original).Error; err != nil {
		return fmt.Errorf("failed archive: %w", err)
	}

	campaigns := []Campaign{}
	if err := db.Model(Campaign{}).
		Where("advertiser_id = ?", original.ID).
		Find(&campaigns).Error; err != nil {

		return fmt.Errorf("error on get camapigns: %w", err)
	}

	for _, c := range campaigns {
		if err := c.Archive(db, archive); err != nil {
			return fmt.Errorf("err on archive campaigns: %w", err)
		}
	}

	return nil
}
