package outputs

import (
	"go.uber.org/fx"

	"go.ads.coffee/server/plugins/outputs"
)

var Module = fx.Module(
	"outputs",

	outputs.Module,

	fx.Provide(
		fx.Annotate(
			New,
			fx.ParamTags(
				`group:"outputs"`,
			),
		),
	),
)
