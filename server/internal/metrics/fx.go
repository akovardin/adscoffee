package metrics

import (
	"go.uber.org/fx"

	"go.ads.coffee/platform/pkg/telemetry"
)

var Module = fx.Module(
	"metrics",
	fx.Invoke(func(tel *telemetry.Telemetry) error {
		return tel.Register(requests, errors, duration)
	}),
)
