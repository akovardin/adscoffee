package media

import (
	"github.com/qor5/admin/v3/media"
	"github.com/qor5/admin/v3/media/oss"
	"github.com/qor5/admin/v3/presets"
	"gorm.io/gorm"

	"go.ads.coffee/platform/admin/internal/s3storage"
)

type Media struct {
	db        *gorm.DB
	s3storage *s3storage.Client
}

func New(db *gorm.DB, s3storage *s3storage.Client) *Media {
	return &Media{
		db:        db,
		s3storage: s3storage,
	}
}

func (m *Media) Configure(b *presets.Builder) {
	oss.Storage = m.s3storage

	mediab := media.New(m.db).AutoMigrate()
	b.Use(mediab)

	// mediab.GetPresetsModelBuilder().
	// 	Label("Медиа")
}
