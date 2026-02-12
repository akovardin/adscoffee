package ads

import (
	"go.uber.org/fx"

	"go.ads.coffee/platform/admin/internal/modules/ads/builders"
)

var Module = fx.Module(
	"ads",
	fx.Provide(
		New,
		builders.NewBanner,
		builders.NewGroup,
		builders.NewCampaign,
		builders.NewAdvertiser,
		builders.NewNetwork,
		builders.NewPlacement,
		builders.NewUnit,
	),
)
