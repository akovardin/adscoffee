package inputs

import "go.uber.org/fx"

var Module = fx.Module(
	"inputs",

	fx.Provide(
		fx.Annotate(
			New,
			fx.ParamTags(
				`group:"inputs"`,
			),
			// fx.ParamTags(`grop:"stages"`),
		),
	),
)
