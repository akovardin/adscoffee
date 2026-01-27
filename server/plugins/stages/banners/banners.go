package banners

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
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

type Banners struct{}

func New() *Banners {
	return &Banners{}
}

func (t *Banners) Name() string {
	return "stages.banners"
}

func (t *Banners) Copy(cfg map[string]any) plugins.Stage {
	return &Banners{}
}

func (t *Banners) Do(ctx context.Context, state *plugins.State) {
	// загружаем баннеры из репозитория
	// и добавляем их в стейт

	state.Candidates = []ads.Banner{}
}
