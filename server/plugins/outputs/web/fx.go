package web

import (
	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/plugins"
	"go.ads.coffee/platform/server/plugins/outputs/web/formats"
)

var Module = fx.Module(
	"outputs.web",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Output)),
			fx.ResultTags(`group:"outputs"`),
			fx.ParamTags(`group:"outputs.web.formats"`),
		),

		formats.NewBanner,
		formats.NewNative,

		fx.Annotate(
			func(b *formats.Banner) plugins.Format {
				return b
			},
			fx.As(new(plugins.Format)),
			fx.ResultTags(`group:"outputs.web.formats"`),
		),

		fx.Annotate(
			func(n *formats.Native) plugins.Format {
				return n
			},
			fx.As(new(plugins.Format)),
			fx.ResultTags(`group:"outputs.web.formats"`),
		),
	),
)
