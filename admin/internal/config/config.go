package config

import (
	"os"

	"go.uber.org/config"
	"go.uber.org/fx"

	"go.ads.coffee/platform/admin/internal/database"
	"go.ads.coffee/platform/admin/internal/s3storage"
	"go.ads.coffee/platform/admin/internal/server"
)

type Config struct {
	fx.Out

	Database  database.Config  `yaml:"database"`
	S3Storage s3storage.Config `yaml:"s3storage"`
	Server    server.Config    `yaml:"server"`
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
