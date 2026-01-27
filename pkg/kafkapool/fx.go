package kafkapool

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"go.ads.coffee/platform/pkg/telemetry"
)

var (
	Module = fx.Module(
		"kafka-pool",
		fx.Provide(
			NewPool,
			NewMetrics,
			adapterTelemetry,
			adapterTracer,
		),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("kafka-pool")
		}),
	)
)

func adapterTelemetry(t *telemetry.Telemetry) Telemetry { //nolint: ireturn
	return t
}

func adapterTracer(t *telemetry.Telemetry) Tracer { //nolint: ireturn
	return t
}
