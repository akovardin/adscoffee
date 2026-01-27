package banners

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
	"go.ads.coffee/platform/server/internal/repos/banners"
)

var Module = fx.Module(
	"stages.banners",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Stage)),
			fx.ResultTags(`group:"stages"`),
		),
	),
)

type BannersCache interface {
	All(ctx context.Context) []ads.Banner
}

type Banners struct {
	cache BannersCache
}

func New(cache *banners.Cache) *Banners {
	return &Banners{
		cache: cache,
	}
}

func (b *Banners) Name() string {
	return "stages.banners"
}

func (b *Banners) Copy(cfg map[string]any) plugins.Stage {
	return &Banners{
		cache: b.cache,
	}
}

func (b *Banners) Do(ctx context.Context, state *plugins.State) error {
	state.Candidates = b.cache.All(ctx)

	return nil
}
