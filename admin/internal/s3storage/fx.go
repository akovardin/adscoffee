package s3storage

import "go.uber.org/fx"

var Module = fx.Module(
	"s3storage",
	fx.Provide(
		New,
	),
)
