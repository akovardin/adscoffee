package stages

import (
	"go.uber.org/fx"

	"go.ads.coffee/server/plugins/stages"
)

var Module = fx.Module(
	"stages",

	stages.Module,

	fx.Provide(
		fx.Annotate(
			New,
			fx.ParamTags(
				`group:"stages"`,
			),
		),
	),
)
