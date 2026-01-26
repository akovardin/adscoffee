package targeting

import (
	"context"

	"go.ads.coffee/platform/server/domain"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"stages.targeting",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(domain.Stage)),
			fx.ResultTags(`group:"stages"`),
		),
	),
)

type Targeting struct{}

func New() *Targeting {
	return &Targeting{}
}

func (t *Targeting) Name() string {
	return "stages.targeting"
}

func (t *Targeting) Copy(cfg map[string]any) domain.Stage {
	return &Targeting{}
}

func (t *Targeting) Targetings(tt []domain.Targeting) {
	// set targetings
}

func (t *Targeting) Do(ctx context.Context, state *domain.State) {
	// обрабатываются таргетинги
	state.Candidates = state.Candidates[:]
}
