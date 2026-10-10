package limits

import (
	"go.ads.coffee/platform/server/plugins/limits"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"limits",

	limits.Module,

	fx.Provide(
		fx.Annotate(
			New,
			fx.ParamTags(
				`group:"limits"`,
			),
		),
	),
)
