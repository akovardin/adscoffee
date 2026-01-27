package config

import (
	"os"

	"go.uber.org/config"
	"go.uber.org/fx"

	"go.ads.coffee/platform/pkg/circuitbreaker"
	"go.ads.coffee/platform/pkg/database"
	"go.ads.coffee/platform/pkg/health"
	"go.ads.coffee/platform/pkg/kafkapool"
	"go.ads.coffee/platform/pkg/redispool"
	"go.ads.coffee/platform/pkg/telemetry"
	"go.ads.coffee/platform/server/internal/pipeline"
)

type Config struct {
	fx.Out

	Pipelines []pipeline.Config `yaml:"pipelines"`

	Health         health.Config                     `yaml:"health"`
	CircuitBreaker map[string]*circuitbreaker.Config `yaml:"circuit-breaker"`
	RedisPool      map[string]*redispool.Config      `yaml:"redis-pool"`
	Telemetry      telemetry.Config                  `yaml:"telemetry"`
	Database       database.Config                   `yaml:"database"`
	Kafka          map[string]*kafkapool.Config      `yaml:"kafka-pool"`
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
