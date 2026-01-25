package inputs

import (
	"go.uber.org/fx"

	"go.ads.coffee/server/plugins/inputs"
)

var Module = fx.Module(
	"inputs",

	inputs.Module,

	fx.Provide(
		fx.Annotate(
			New,
			fx.ParamTags(
				`group:"inputs"`,
			),
		),
	),
)
