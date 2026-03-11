package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
	"go.uber.org/fx"

	"go.ads.coffee/platform/analytics/internal/clickhouse"
	"go.ads.coffee/platform/analytics/internal/config"
	"go.ads.coffee/platform/analytics/internal/handlers"
)

func main() {
	cmd := &cli.Command{
		Name: "analytics",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "config", Aliases: []string{"c"}},
		},
		Commands: []*cli.Command{
			{
				Name: "aggregations",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "table", Aliases: []string{"t"}},
				},
				Aliases: []string{"a"},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fx.New(
						fx.Provide(
							func() (config.Config, error) {
								cfg := cmd.String("config")
								if cfg == "" {
									cfg = "analytics/configs/config.yaml"
								}

								return config.New(cfg)
							},
						),

						handlers.Module,
						clickhouse.Module,

						fx.Invoke(
							func(aggregate *handlers.Aggregate) {
								if err := aggregate.Run(ctx, cmd); err != nil {
									panic(err)
								}

								os.Exit(0)
							},
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
}
