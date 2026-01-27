package tracker

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"inputs.tracker",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Input)),
			fx.ResultTags(`group:"inputs"`),
		),
	),
)

type Tracker struct {
}

func New() *Tracker {
	return &Tracker{}
}

func (s *Tracker) Name() string {
	return "inputs.tracker"
}

func (s *Tracker) Copy(cfg map[string]any) plugins.Input {
	return &Tracker{}
}

func (stages *Tracker) Do(ctx context.Context, state *plugins.State) bool {
	// нужно получить данные пользователя из запроса

	state.User = &plugins.User{}
	state.Device = &plugins.Device{}

	return true
}
