package units

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"repos.units",

	fx.Provide(
		NewRepo,
		NewCache,
	),
)
