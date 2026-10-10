package inapp

import (
	"go.uber.org/fx"

	"go.ads.coffee/platform/server/plugins/outputs/inapp/formats"
)

var Module = fx.Module(
	"outputs.inapp",

	fx.Provide(
		fx.Annotate(
			New,
			fx.ResultTags(`group:"outputs"`),
			fx.ParamTags(`group:"outputs.inapp.formats"`),
		),

		fx.Annotate(
			formats.NewNative,
			fx.ResultTags(`group:"outputs.inapp.formats"`),
		),
	),
)
