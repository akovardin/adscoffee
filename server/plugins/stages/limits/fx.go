package limits

import (
	"go.uber.org/fx"

	"go.ads.coffee/server/domain"
)

var Module = fx.Module(
	"stages.limits",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(domain.Stage)),
			fx.ResultTags(`group:"stages"`),
		),
	),
)
