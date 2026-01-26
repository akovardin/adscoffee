package banners

import (
	"context"

	"go.ads.coffee/platform/server/domain"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"stages.banners",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(domain.Stage)),
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

func (t *Banners) Copy(cfg map[string]any) domain.Stage {
	return &Banners{}
}

func (t *Banners) Do(ctx context.Context, state *domain.State) {
	// загружаем баннеры из репозитория
	// и добавляем их в стейт

	state.Candidates = []domain.Banner{}
}
