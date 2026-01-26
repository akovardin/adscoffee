package tracker

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/config"
	"go.ads.coffee/platform/server/domain"
)

var Module = fx.Module(
	"inputs.tracker",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(domain.Input)),
			fx.ResultTags(`group:"inputs"`),
		),
	),
)

type Tracker struct {
}

func New(config config.Config) *Tracker {
	return &Tracker{}
}

func (s *Tracker) Name() string {
	return "inputs.tracker"
}

func (s *Tracker) Copy(cfg map[string]any) domain.Input {
	return &Tracker{}
}

func (stages *Tracker) Do(ctx context.Context, state *domain.State) bool {
	// нужно получить данные пользователя из запроса

	state.User = &domain.User{}
	state.Device = &domain.Device{}

	return true
}
