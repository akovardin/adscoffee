package builders

import (
	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go.ads.coffee/platform/admin/internal/modules/ads/models"
)

type Unit struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewUnit(logger *zap.Logger, db *gorm.DB) *Unit {
	return &Unit{
		logger: logger,
		db:     db,
	}
}

func (n *Unit) Configure(b *presets.Builder) {
	mn := b.Model(&models.Unit{}).
		MenuIcon("mdi-lan").
		// Label("Рекламодатели").
		RightDrawerWidth("1000")

	mn.Listing("ID", "Name")

	mn.Editing().ValidateFunc(func(obj interface{}, ctx *web.EventContext) (err web.ValidationErrors) {
		u := obj.(*models.Unit)

		if u.Name == "" {
			err.FieldError("Name", "Name is required")
		}
		return
	})
}
