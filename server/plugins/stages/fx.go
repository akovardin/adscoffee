package stages

import "go.uber.org/fx"

var Module = fx.Module(
	"stages",

	fx.Provide(
		fx.Annotate(
			New,
			fx.ParamTags(
				`group:"stages"`,
			),
		),
	),
)
