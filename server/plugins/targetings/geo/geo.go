package geo

import (
	"go.ads.coffee/server/domain"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"targetings.geo",
	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(domain.Targeting)),
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

func (g *Geo) Copy(cfg map[string]any) domain.Targeting {
	return &Geo{}
}

func (g *Geo) Filter() {

}
