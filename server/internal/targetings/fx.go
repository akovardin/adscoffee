package targetings

import (
	"go.ads.coffee/server/plugins/targetings"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"targetings",

	targetings.Module,

	fx.Provide(
		fx.Annotate(
			New,
			fx.ParamTags(
				`group:"targetings"`,
			),
		),
	),
)
