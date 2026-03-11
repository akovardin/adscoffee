package config

import (
	"os"

	"go.uber.org/config"
	"go.uber.org/fx"

	"go.ads.coffee/platform/analytics/internal/clickhouse"
)

type Config struct {
	fx.Out

	Clickhouse clickhouse.Config `yaml:"clickhouse"`
}

func New(file string) (Config, error) {
	provider, err := config.NewYAML(
		config.Expand(os.LookupEnv),
		config.File(file),
		config.Permissive(),
	)

	if err != nil {
		return Config{}, err
	}

	cfg := Config{}

	err = provider.Get("").Populate(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
