package users

import "go.uber.org/fx"

var Module = fx.Module(
	"users",
	fx.Provide(
		New,
	),
)
