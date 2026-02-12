package builders

import (
	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go.ads.coffee/platform/admin/internal/modules/ads/models"
)

type Placement struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewPlacement(logger *zap.Logger, db *gorm.DB) *Placement {
	return &Placement{
		logger: logger,
		db:     db,
	}
}

func (n *Placement) Configure(b *presets.Builder) {
	mn := b.Model(&models.Placement{}).
		MenuIcon("mdi-lan").
		// Label("Рекламодатели").
		RightDrawerWidth("1000")

	mn.Listing("ID", "Name")

	mn.Editing().ValidateFunc(func(obj interface{}, ctx *web.EventContext) (err web.ValidationErrors) {
		u := obj.(*models.Placement)

		if u.Name == "" {
			err.FieldError("Name", "Name is required")
		}
		return
	})
}
