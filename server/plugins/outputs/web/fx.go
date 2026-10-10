package web

import (
	"go.uber.org/fx"

	"go.ads.coffee/platform/server/plugins/outputs/web/formats"
)

var Module = fx.Module(
	"outputs.web",

	fx.Provide(
		fx.Annotate(
			New,
			fx.ResultTags(`group:"outputs"`),
			fx.ParamTags(`group:"outputs.web.formats"`),
		),

		fx.Annotate(
			formats.NewNative,
			fx.ResultTags(`group:"outputs.web.formats"`),
		),
	),
)
