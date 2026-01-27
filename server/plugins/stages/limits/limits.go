package limits

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"stages.limits",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Stage)),
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

func (l *Limits) Copy(cfg map[string]any) plugins.Stage {
	return &Limits{}
}

func (l *Limits) Do(ctx context.Context, state *plugins.State) error {
	// срабатывают ограничения по показам, капингам и так далее
	state.Candidates = state.Candidates[:]

	return nil
}
