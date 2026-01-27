package clock

import "go.uber.org/fx"

var Module = fx.Module(
	"clock",
	fx.Provide(
		New,
	),
)
