package targeting

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"stages.targeting",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Stage)),
			fx.ResultTags(`group:"stages"`),
		),
	),
)

type Targeting struct {
	targetings []plugins.Targeting
}

func New() *Targeting {
	return &Targeting{}
}

func (t *Targeting) Name() string {
	return "stages.targeting"
}

func (t *Targeting) Copy(cfg map[string]any) plugins.Stage {
	return &Targeting{
		targetings: t.targetings,
	}
}

func (t *Targeting) Targetings(tt []plugins.Targeting) {
	t.targetings = tt
}

func (t *Targeting) Do(ctx context.Context, state *plugins.State) error {
	// последовательно применяем таргетинги к кандидатам
	for _, targeting := range t.targetings {
		state.Candidates = targeting.Filter(state.Candidates, state)
	}

	return nil
}
