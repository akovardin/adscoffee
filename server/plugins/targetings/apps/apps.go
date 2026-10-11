package apps

import (
	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"targetings.apps",
	fx.Provide(
		fx.Annotate(
			New,
			fx.ResultTags(`group:"targetings"`),
		),
	),
)

type Apps struct{}

func New() plugins.Provider[plugins.Targeting] {
	return plugins.Provider[plugins.Targeting]{
		Name: "targetings.apps",
		Factory: func(cfg map[string]any) (plugins.Targeting, error) {
			return &Apps{}, nil
		},
	}
}

func (a *Apps) Name() string {
	return "targetings.apps"
}

func (a *Apps) Filter(candidates []*ads.Banner, state *plugins.State) []*ads.Banner {
	return candidates
}
