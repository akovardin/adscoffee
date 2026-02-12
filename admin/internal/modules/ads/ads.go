package ads

import (
	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go.ads.coffee/platform/admin/internal/modules/ads/builders"
	"go.ads.coffee/platform/admin/internal/modules/ads/models"
)

// конфигурация админки для модуля ads
type Ads struct {
	logger *zap.Logger
	db     *gorm.DB

	banner     *builders.Banner
	group      *builders.Group
	campaign   *builders.Campaign
	advertiser *builders.Advertiser
	network    *builders.Network
	placement  *builders.Placement
	unit       *builders.Unit
}

func New(
	logger *zap.Logger,
	db *gorm.DB,
	banner *builders.Banner,
	group *builders.Group,
	campaign *builders.Campaign,
	advertiser *builders.Advertiser,
	network *builders.Network,
	placement *builders.Placement,
	unit *builders.Unit,
) *Ads {
	return &Ads{
		logger:     logger,
		db:         db,
		banner:     banner,
		group:      group,
		campaign:   campaign,
		advertiser: advertiser,
		network:    network,
		placement:  placement,
		unit:       unit,
	}
}

func (m *Ads) Configure(b *presets.Builder) {
	b.AssetFunc(func(ctx *web.EventContext) {
		ctx.Injector.HeadHTML(`
			<style>
			details summary::marker {
				content: "❯ ";
			}
			
			/*
			details[open] summary::marker { content:" " }
			*/
			</style> 
		`)
	})

	m.advertiser.Configure(b)
	m.campaign.Configure(b)
	m.group.Configure(b)
	m.banner.Configure(b)
	m.network.Configure(b)
	m.placement.Configure(b)
	m.unit.Configure(b)
}

// TODO mowe to different command
func (u *Ads) Migrate() {
	err := u.db.AutoMigrate(
		&models.Advertiser{},
		&models.Campaign{},
		&models.Banner{},
		&models.Bgroup{},
		&models.Audience{},
		&models.Network{},
	)
	if err != nil {
		panic(err)
	}
}
