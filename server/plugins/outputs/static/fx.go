package static

import (
	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/plugins"
	"go.ads.coffee/platform/server/plugins/outputs/static/formats"
)

var Module = fx.Module(
	"outputs.static",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Output)),
			fx.ResultTags(`group:"outputs"`),
			fx.ParamTags(`group:"outputs.static.formats"`),
		),

		formats.NewBanner,

		fx.Annotate(
			func(b *formats.Banner) plugins.Format {
				return b
			},
			fx.As(new(plugins.Format)),
			fx.ResultTags(`group:"outputs.static.formats"`),
		),
	),
)
