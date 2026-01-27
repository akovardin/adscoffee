package rotation

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"stages.mediation",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Stage)),
			fx.ResultTags(`group:"stages"`),
		),
	),
)

type Mediation struct{}

func New() *Mediation {
	return &Mediation{}
}

func (t *Mediation) Name() string {
	return "stages.mediation"
}

func (t *Mediation) Copy(cfg map[string]any) plugins.Stage {
	return &Mediation{}
}

func (t *Mediation) Do(ctx context.Context, state *plugins.State) {
	// делаем взвешанный рандом по ecpm
	state.Winners = state.Candidates

	// тут нужно взять фйековые баннера для в зависимости от настройки
	// медиации на юните

}
