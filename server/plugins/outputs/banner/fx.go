package banner

import (
	"go.uber.org/fx"

	"go.ads.coffee/server/domain"
)

var Module = fx.Module(
	"banner",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(domain.Output)),
			fx.ResultTags(`group:"outputs"`),
		),
	),
)
