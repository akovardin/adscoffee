package stats

import "go.uber.org/fx"

var Module = fx.Module(
	"stats",
	fx.Provide(
		New,
		NewQuery,
	),
)
