package inapp

import (
	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/plugins"
	"go.ads.coffee/platform/server/plugins/outputs/inapp/formats"
)

var Module = fx.Module(
	"outputs.inapp",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Output)),
			fx.ResultTags(`group:"outputs"`),
			fx.ParamTags(`group:"outputs.inapp.formats"`),
		),

		formats.NewNative,

		fx.Annotate(
			func(n *formats.Native) plugins.Format {
				return n
			},
			fx.As(new(plugins.Format)),
			fx.ResultTags(`group:"outputs.inapp.formats"`),
		),
	),
)