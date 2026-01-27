package geo

import (
	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"targetings.geo",
	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Targeting)),
			fx.ResultTags(`group:"targetings"`),
		),
	),
)

type Geo struct {
}

func New() *Geo {
	return &Geo{}
}

func (g *Geo) Name() string {
	return "targetings.geo"
}

func (g *Geo) Copy(cfg map[string]any) plugins.Targeting {
	return &Geo{}
}

func (g *Geo) Filter() {

}
