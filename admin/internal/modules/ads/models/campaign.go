package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Campaign struct {
	gorm.Model

	Title  string
	Active bool

	Bundle string

	Start time.Time
	End   time.Time

	Targeting string
	Budget    string
	Capping   string
	Timetable string

	AdvertiserID int
	Advertiser   Advertiser

	ArchivedAt *time.Time
}

func (original Campaign) Copy(db *gorm.DB, advertiser int) (Campaign, error) {
	copy := Campaign{
		AdvertiserID: advertiser,
		Title:        original.Title + " (Копия)",
		Bundle:       original.Bundle,
		Start:        original.Start,
		End:          original.End,
		Timetable:    original.Timetable,
		Targeting:    original.Targeting,
		Budget:       original.Budget,
		Capping:      original.Capping,
		Active:       false,
	}

	// Сохраняем копию в базу данных
	if err := db.Create(&copy).Error; err != nil {
		return Campaign{}, fmt.Errorf("failed to create copy: %w", err)
	}

	groups := []Bgroup{}
	if err := db.Model(Bgroup{}).
		Where("campaign_id = ?", original.ID).
		Find(&groups).Error; err != nil {
		return Campaign{}, fmt.Errorf("error on get groups: %w", err)
	}

	for _, g := range groups {
		if _, err := g.Copy(db, int(copy.ID)); err != nil {
			return Campaign{}, fmt.Errorf("err on copy groups: %w", err)
		}
	}

	return copy, nil
}

func (original Campaign) Archive(db *gorm.DB, archive *time.Time) error {
	original.ArchivedAt = archive

	if err := db.Save(original).Error; err != nil {
		return fmt.Errorf("failed archive: %w", err)
	}

	groups := []Bgroup{}
	if err := db.Model(Bgroup{}).
		Where("campaign_id = ?", original.ID).
		Find(&groups).Error; err != nil {

		return fmt.Errorf("error on get groups: %w", err)
	}

	for _, g := range groups {
		if err := g.Archive(db, archive); err != nil {
			return fmt.Errorf("err on archive groups: %w", err)
		}
	}

	return nil
}
