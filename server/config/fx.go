package config

import (
	"flag"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"config",
	fx.Provide(
		func() (Config, error) {
			var (
				cfg string
			)

			flag.StringVar(&cfg, "c", "server/config.yaml", "config")
			flag.Parse()

			return New(cfg)
		},
	),
)
