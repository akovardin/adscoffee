package budget

import (
	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"limits.budget",
	fx.Provide(
		fx.Annotate(
			New,
			fx.ResultTags(`group:"limits"`),
		),
	),
)

// Budget — лимит по бюджету рекламодателя. Пока заглушка.
type Budget struct{}

func New() plugins.Provider[plugins.Limit] {
	return plugins.Provider[plugins.Limit]{
		Name: "limits.budget",
		Factory: func(cfg map[string]any) (plugins.Limit, error) {
			return &Budget{}, nil
		},
	}
}

func (b *Budget) Name() string {
	return "limits.budget"
}

func (b *Budget) Apply(candidates []*ads.Banner, state *plugins.State) []*ads.Banner {
	state.Candidates = state.Candidates[:]

	return candidates
}
