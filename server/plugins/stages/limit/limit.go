package limit

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"stages.limit",

	fx.Provide(
		fx.Annotate(
			New,
			fx.ResultTags(`group:"stages"`),
		),
	),
)

type Limit struct {
	limits []plugins.Limit
}

func New() plugins.Provider[plugins.Stage] {
	return plugins.Provider[plugins.Stage]{
		Name: "stages.limit",
		Factory: func(cfg map[string]any) (plugins.Stage, error) {
			return &Limit{}, nil
		},
	}
}

func (l *Limit) Name() string {
	return "stages.limit"
}

func (l *Limit) Limits(ll []plugins.Limit) {
	l.limits = ll
}

func (l *Limit) Do(ctx context.Context, state *plugins.State) error {
	// последовательно применяем лимиты к кандидатам
	for _, limit := range l.limits {
		state.Candidates = limit.Apply(state.Candidates, state)
	}

	return nil
}
