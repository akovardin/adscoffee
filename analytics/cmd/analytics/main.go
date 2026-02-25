package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
	"go.ads.coffee/platform/analytics/internal/handlers"
	"go.uber.org/fx"
)

func main() {
	cmd := &cli.Command{
		Name: "analytics",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "config", Aliases: []string{"c"}},
		},
		Commands: []*cli.Command{
			{
				Name:    "impressions",
				Aliases: []string{"i"},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fx.New(
						fx.Provide(
							// func() (config.Config, error) {
							// 	cfg := cmd.String("config")
							// 	if cfg == "" {
							// 		cfg = "analytics/configs/config.yaml"
							// 	}

							// 	return config.New(cfg)
							// },

							handlers.Module,
						),

						fx.Invoke(
							func(impressions *handlers.Impressions) {
								impressions.Run()
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
