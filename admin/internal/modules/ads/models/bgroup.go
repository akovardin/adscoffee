package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Bgroup struct {
	gorm.Model

	Title  string
	Active bool

	Price int

	Start time.Time
	End   time.Time

	Targeting string
	Budget    string
	Capping   string
	Timetable string

	CampaignID int
	Campaign   Campaign

	ArchivedAt *time.Time
}

func (original Bgroup) Copy(db *gorm.DB, campaign int) (Bgroup, error) {
	copy := Bgroup{
		CampaignID: campaign,
		Title:      original.Title + " (Копия)",
		Price:      original.Price,
		Start:      original.Start,
		End:        original.End,
		Timetable:  original.Timetable,
		Targeting:  original.Targeting,
		Budget:     original.Budget,
		Capping:    original.Capping,
		Active:     false,
	}

	if err := db.Create(&copy).Error; err != nil {
		return Bgroup{}, fmt.Errorf("failed to create copy: %w", err)
	}

	banners := []Banner{}
	if err := db.Model(Banner{}).
		Where("bgroup_id = ?", original.ID).
		Find(&banners).Error; err != nil {
		return Bgroup{}, fmt.Errorf("error on get banners: %w", err)
	}

	for _, b := range banners {
		if _, err := b.Copy(db, int(copy.ID)); err != nil {
			return Bgroup{}, fmt.Errorf("err on copy banner: %w", err)
		}
	}

	return copy, nil
}

func (original Bgroup) Archive(db *gorm.DB, archive *time.Time) error {
	original.ArchivedAt = archive

	if err := db.Save(&original).Error; err != nil {
		return fmt.Errorf("failed to archive: %w", err)
	}

	banners := []Banner{}
	if err := db.Model(Banner{}).
		Where("bgroup_id = ?", original.ID).
		Find(&banners).Error; err != nil {
		return fmt.Errorf("error on get banners: %w", err)
	}

	for _, b := range banners {
		if err := b.Archive(db, archive); err != nil {
			return fmt.Errorf("err on archive banner: %w", err)
		}
	}

	return nil
}
