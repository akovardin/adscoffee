package builders

import (
	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go.ads.coffee/platform/admin/internal/modules/ads/models"
)

type Network struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewNetwork(logger *zap.Logger, db *gorm.DB) *Network {
	return &Network{
		logger: logger,
		db:     db,
	}
}

func (n *Network) Configure(b *presets.Builder) {
	mn := b.Model(&models.Network{}).
		MenuIcon("mdi-lan").
		// Label("Рекламодатели").
		RightDrawerWidth("1000")

	mn.Listing("ID", "Title", "Name")

	mn.Editing().ValidateFunc(func(obj interface{}, ctx *web.EventContext) (err web.ValidationErrors) {
		u := obj.(*models.Network)
		if u.Title == "" {
			err.FieldError("Title", "Title is required")
		}
		if u.Name == "" {
			err.FieldError("Name", "Name is required")
		}
		return
	})
}
