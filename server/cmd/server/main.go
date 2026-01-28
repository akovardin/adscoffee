package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/fx"

	"go.ads.coffee/platform/pkg/database"
	"go.ads.coffee/platform/pkg/logger"
	"go.ads.coffee/platform/server/internal/config"
	"go.ads.coffee/platform/server/internal/repos/banners"
	"go.ads.coffee/platform/server/internal/repos/placements"
	"go.ads.coffee/platform/server/internal/server"
	"go.ads.coffee/platform/server/plugins"
)

func main() {
	app := fx.New(
		config.Module,
		logger.Module,
		server.Module,
		database.Module,
		plugins.Module,

		// repos
		banners.Module,
		placements.Module,

		fx.Invoke(
			start,
			caches,
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

func caches(banners *banners.Cache) {
	go banners.Start(context.Background())
}
