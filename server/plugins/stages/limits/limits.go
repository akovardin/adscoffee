package limits

import (
	"context"

	"go.ads.coffee/platform/server/domain"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"stages.limits",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(domain.Stage)),
			fx.ResultTags(`group:"stages"`),
		),
	),
)

type Limits struct{}

func New() *Limits {
	return &Limits{}
}

func (l *Limits) Name() string {
	return "stages.limits"
}

func (l *Limits) Copy(cfg map[string]any) domain.Stage {
	return &Limits{}
}

func (l *Limits) Do(ctx context.Context, state *domain.State) {
	// срабатывают ограничения по показам, капингам и так далее
	state.Candidates = state.Candidates[:]
}
