package main

import (
	"context"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"

	"go.ads.coffee/platform/pkg/circuitbreaker"
	"go.ads.coffee/platform/pkg/database"
	"go.ads.coffee/platform/pkg/health"
	"go.ads.coffee/platform/pkg/kafkapool"
	"go.ads.coffee/platform/pkg/logger"
	"go.ads.coffee/platform/pkg/redispool"
	"go.ads.coffee/platform/pkg/telemetry"
	"go.ads.coffee/platform/server/internal/analytics"
	"go.ads.coffee/platform/server/internal/config"
	"go.ads.coffee/platform/server/internal/repos/banners"
	"go.ads.coffee/platform/server/internal/repos/placements"
	"go.ads.coffee/platform/server/internal/server"
	"go.ads.coffee/platform/server/internal/sessions"
	"go.ads.coffee/platform/server/plugins"
)

func main() {
	cmd := &cli.Command{
		Name: "kodikapusta",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "config", Aliases: []string{"c"}},
		},
		Commands: []*cli.Command{
			{
				Name:    "serve",
				Aliases: []string{"s"},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fx.New(
						fx.Provide(
							func() prometheus.Registerer {
								// default prometheus
								return prometheus.DefaultRegisterer
							},
						),
						fx.Provide(
							func() (config.Config, error) {
								cfg := cmd.String("config")
								if cfg == "" {
									cfg = "server/configs/config.yaml"
								}

								return config.New(cfg)
							},
						),
						logger.Module,
						server.Module,
						database.Module,
						sessions.Module,
						analytics.Module,
						telemetry.Module,
						health.Module,
						circuitbreaker.Module,
						redispool.Module,
						kafkapool.Module,
						plugins.Module,

						// repos
						banners.Module,
						placements.Module,

						fx.Invoke(
							start,
							caches,
						),
					).Run()

					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		panic(err)
	}

	// sig := make(chan os.Signal, 1)
	// signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	// <-sig

	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()

	// if err := app.Stop(ctx); err != nil {
	// 	panic(err)
	// }
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

func caches(banners *banners.Cache) {
	go banners.Start(context.Background())
}
