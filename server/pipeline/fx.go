package pipeline

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"pipeline",

	fx.Provide(
		NewManager,
	),
)
