package circuitbreaker

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"go.ads.coffee/platform/pkg/telemetry"
)

var Module = fx.Module(
	"circuitbreaker",
	fx.Provide(
		NewPool,
		NewMetrics,
		adapterTelemetry,
	),
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("circuitbreaker")
	}),
)

func adapterTelemetry(t *telemetry.Telemetry) Telemetry { //nolint: ireturn
	return t
}
