package sessions

import "go.uber.org/fx"

var Module = fx.Module(
	"sessions",
	fx.Provide(
		New,
	),
)
