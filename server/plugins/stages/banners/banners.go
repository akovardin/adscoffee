package banners

import (
	"context"
	"slices"

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
			fx.ResultTags(`group:"stages"`),
		),
	),
)

type BannersCache interface {
	All(ctx context.Context) []*ads.Banner
}

type Banners struct {
	cache BannersCache
}

func New(cache *banners.Cache) plugins.Provider[plugins.Stage] {
	return plugins.Provider[plugins.Stage]{
		Name: "stages.banners",
		Factory: func(cfg map[string]any) (plugins.Stage, error) {
			return &Banners{cache: cache}, nil
		},
	}
}

func (b *Banners) Name() string {
	return "stages.banners"
}

func (b *Banners) Do(ctx context.Context, state *plugins.State) error {
	// копируем только слайс указателей (n*8 байт): дальше таргетинги/лимиты
	// фильтруют in-place, не трогая общий кэш баннеров
	state.Candidates = slices.Clone(b.cache.All(ctx))

	return nil
}
