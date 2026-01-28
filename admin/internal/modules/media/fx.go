package media

import "go.uber.org/fx"

var Module = fx.Module(
	"media",
	fx.Provide(
		New,
	),
)
