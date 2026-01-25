package geo

import (
	"go.uber.org/fx"

	"go.ads.coffee/server/domain"
)

var Module = fx.Module(
	"targetings.geo",
	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(domain.Targeting)),
			fx.ResultTags(`group:"targetings"`),
		),
	),
)
