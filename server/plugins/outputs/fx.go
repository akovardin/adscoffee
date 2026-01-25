package outputs

import "go.uber.org/fx"

var Module = fx.Module(
	"outputs",

	fx.Provide(
		fx.Annotate(
			New,
			fx.ParamTags(
				`group:"outputs"`,
			),
		),
	),
)
