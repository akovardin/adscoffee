package redispool

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"go.ads.coffee/platform/pkg/health"
	"go.ads.coffee/platform/pkg/telemetry"
)

var Module = fx.Module(
	"redispool",
	fx.Provide(
		NewPool,
		NewMetrics,
		adapterHealth,
		adapterTelemetry,
		adapterTracer,
	),
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("redispool")
	}),
	fx.Invoke(
		func(lc fx.Lifecycle, pp *Pool) {
			lc.Append(fx.Hook{
				OnStop: pp.Stop,
			})
		},
	),
)

type HealthComponentOut struct {
	fx.Out

	Pool *health.Component `group:"health_component"`
}

func adapterHealth(pp *Pool) HealthComponentOut {
	return HealthComponentOut{
		Pool: pp.HealthComponent(),
	}
}

func adapterTelemetry(t *telemetry.Telemetry) Telemetry { //nolint: ireturn
	return t
}

func adapterTracer(t *telemetry.Telemetry) Tracer { //nolint: ireturn
	return t
}
