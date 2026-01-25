package targetings

import "go.uber.org/fx"

var Module = fx.Module(
	"targetings",

	fx.Provide(
		fx.Annotate(
			New,
			fx.ParamTags(
				`group:"targetings"`,
			),
		),
	),
)
