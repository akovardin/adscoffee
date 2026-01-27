package health

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"health",
	fx.Provide(
		NewHealth,
	),
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("health")
	}),
)

type Params struct {
	fx.In

	Logger *zap.Logger
	Config Config

	Components         []*Component        `group:"health_component"`
	ComponentProviders []ComponentProvider `group:"health_component_provider"`
}
