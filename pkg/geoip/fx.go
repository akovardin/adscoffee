package geoip

import "go.uber.org/fx"

var Module = fx.Module(
	"geoip",
	fx.Provide(
		New,
	),
)
