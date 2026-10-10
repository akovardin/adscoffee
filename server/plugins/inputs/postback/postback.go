package postback

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/plugins"
	"go.ads.coffee/platform/server/internal/metrics"
)

var Module = fx.Module(
	"inputs.postback",

	fx.Provide(
		fx.Annotate(
			New,
			fx.ResultTags(`group:"inputs"`),
		),
	),
)

type Postback struct {
}

func New() plugins.Provider[plugins.Input] {
	return plugins.Provider[plugins.Input]{
		Name: "inputs.postback",
		Factory: func(cfg map[string]any) (plugins.Input, error) {
			return &Postback{}, nil
		},
	}
}

func (s *Postback) Name() string {
	return "inputs.postback"
}

func (stages *Postback) Do(ctx context.Context, state *plugins.State) bool {
	metrics.IncInput(stages.Name())

	// нужно получить данные пользователя из запроса

	state.User = &plugins.User{}
	state.Device = &plugins.Device{}

	return true
}
