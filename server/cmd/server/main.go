package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/config"
	"go.uber.org/fx"

	cfg "go.ads.coffee/platform/server/config"
	"go.ads.coffee/platform/server/internal/server"
	"go.ads.coffee/platform/server/plugins"
)

func main() {
	app := fx.New(
		fx.Provide(
			func() (cfg.Config, error) {
				provider, err := config.NewYAML(
					config.Expand(os.LookupEnv),
					config.File("server/config.yaml"),
					config.Permissive(),
				)
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
	)

	app.Run()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Stop(ctx); err != nil {
		panic(err)
	}
}

func start(lc fx.Lifecycle, server *server.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return server.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}
