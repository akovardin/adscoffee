package rtb

import (
	"go.uber.org/fx"

	"go.ads.coffee/server/domain"
)

var Module = fx.Module(
	"rtb",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(domain.Input)),
			fx.ResultTags(`group:"inputs"`),
		),
	),
)
