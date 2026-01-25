package main

import (
	"context"

	"go.uber.org/config"
	"go.uber.org/fx"

	cfg "go.ads.coffee/server/config"
	"go.ads.coffee/server/internal/server"
	"go.ads.coffee/server/plugins"
)

func main() {
	fx.New(
		fx.Provide(
			func() (cfg.Config, error) {

				base := config.File("/Users/artem/projects/adscoffee/platform/server/config.yaml")

				provider, err := config.NewYAML(base)
				if err != nil {
					return cfg.Config{}, err
				}

				c := cfg.Config{}

				if err := provider.Get("").Populate(&c); err != nil {
					return cfg.Config{}, err
				}

				return c, nil
			},
		),

		plugins.Module,
		server.Module,

		fx.Invoke(
			start,
		),
	).Run()
}

func start(server *server.Server) {
	server.Start(context.Background())
}
