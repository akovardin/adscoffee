package models

import (
	"fmt"
	"time"

	"github.com/qor5/admin/v3/media/media_library"
	"gorm.io/gorm"
)

type Banner struct {
	gorm.Model

	Title       string
	Label       string
	Description string
	Active      bool

	Erid         string
	OrdCategory  string
	OrdTargeting string
	OrdFormat    string
	OrdKktu      string

	Price int

	Image media_library.MediaBox `sql:"type:text;"`
	Icon  media_library.MediaBox `sql:"type:text;"`

	Start time.Time
	End   time.Time

	Clicktracker string
	Imptracker   string
	Target       string

	Targeting string
	Budget    string
	Capping   string

	BgroupID int
	Bgroup   Bgroup

	Timetable       string
	ExpectedWinRate float64
	ArchivedAt      *time.Time
}

func (original Banner) Copy(db *gorm.DB, group int) (Banner, error) {
	copy := Banner{
		BgroupID:    group,
		Title:       original.Title + " (Копия)",
		Label:       original.Label,
		Description: original.Description,
		Start:       original.Start,
		End:         original.End,
		Timetable:   original.Timetable,
		Targeting:   original.Targeting,
		Budget:      original.Budget,
		Capping:     original.Capping,
		Price:       original.Price,
		Active:      false,

		Erid:         original.Erid,
		OrdCategory:  original.OrdCategory,
		OrdTargeting: original.OrdTargeting,
		OrdFormat:    original.OrdFormat,
		OrdKktu:      original.OrdKktu,

		Image: original.Image,
		Icon:  original.Icon,

		Clicktracker: original.Clicktracker,
		Imptracker:   original.Imptracker,
		Target:       original.Target,
	}

	if err := db.Create(&copy).Error; err != nil {
		return copy, fmt.Errorf("failed to create copy banners: %w", err)
	}

	return copy, nil
}
