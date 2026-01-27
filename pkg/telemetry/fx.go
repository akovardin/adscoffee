package telemetry

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"telemetry",
	fx.Provide(
		New,
	),
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("telemetry")
	}),
	fx.Invoke(func(lc fx.Lifecycle, tel *Telemetry) {
		lc.Append(fx.Hook{
			OnStop: tel.Stop,
		})
	}),
)
