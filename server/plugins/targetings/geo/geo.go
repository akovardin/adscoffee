package geo

import (
	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"targetings.geo",
	fx.Provide(
		fx.Annotate(
			New,
			fx.ResultTags(`group:"targetings"`),
		),
	),
)

type Geo struct {
}

func New() plugins.Provider[plugins.Targeting] {
	return plugins.Provider[plugins.Targeting]{
		Name: "targetings.geo",
		Factory: func(cfg map[string]any) (plugins.Targeting, error) {
			return &Geo{}, nil
		},
	}
}

func (g *Geo) Name() string {
	return "targetings.geo"
}

func (g *Geo) Filter(candidates []*ads.Banner, state *plugins.State) []*ads.Banner {
	return candidates
}
