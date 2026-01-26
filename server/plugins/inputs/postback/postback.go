package postback

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/config"
	"go.ads.coffee/platform/server/domain"
)

var Module = fx.Module(
	"inputs.postback",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(domain.Input)),
			fx.ResultTags(`group:"inputs"`),
		),
	),
)

type Postback struct {
}

func New(config config.Config) *Postback {
	return &Postback{}
}

func (s *Postback) Name() string {
	return "inputs.postback"
}

func (s *Postback) Copy(cfg map[string]any) domain.Input {
	return &Postback{}
}

func (stages *Postback) Do(ctx context.Context, state *domain.State) bool {
	// нужно получить данные пользователя из запроса

	state.User = &domain.User{}
	state.Device = &domain.Device{}

	return true
}
