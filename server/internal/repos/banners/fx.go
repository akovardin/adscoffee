package banners

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"repos.banner",

	fx.Provide(
		NewRepo,
		NewCache,
	),
)
