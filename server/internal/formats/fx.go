package formats

import (
	"go.uber.org/fx"

	"go.ads.coffee/platform/server/plugins/formats"
)

var Module = fx.Module(
	"formats",

	formats.Module,

	fx.Provide(
		fx.Annotate(
			New,
			fx.ParamTags(
				`group:"formats"`,
			),
		),
	),
)
