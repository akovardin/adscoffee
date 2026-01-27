package placements

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"repos.placements",

	fx.Provide(
		NewRepo,
		NewCache,
	),
)
