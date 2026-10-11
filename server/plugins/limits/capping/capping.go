package capping

import (
	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"limits.capping",
	fx.Provide(
		fx.Annotate(
			New,
			fx.ResultTags(`group:"limits"`),
		),
	),
)

// Capping — лимит по показам/кликам (капинг). Пока заглушка.
type Capping struct{}

func New() plugins.Provider[plugins.Limit] {
	return plugins.Provider[plugins.Limit]{
		Name: "limits.capping",
		Factory: func(cfg map[string]any) (plugins.Limit, error) {
			return &Capping{}, nil
		},
	}
}

func (c *Capping) Name() string {
	return "limits.capping"
}

func (c *Capping) Apply(candidates []*ads.Banner, state *plugins.State) []*ads.Banner {
	state.Candidates = state.Candidates[:]

	return candidates
}
