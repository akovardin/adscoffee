package config

import (
	"os"

	"go.ads.coffee/platform/pkg/database"
	"go.ads.coffee/platform/server/internal/pipeline"
	"go.uber.org/config"
	"go.uber.org/fx"
)

type Config struct {
	fx.Out

	Pipelines []pipeline.Config `yaml:"pipelines"`
	Database  database.Config   `yaml:"database"`
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
