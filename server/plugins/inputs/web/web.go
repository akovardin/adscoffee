package web

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/config"
	"go.ads.coffee/platform/server/domain"
)

var Module = fx.Module(
	"inputs.web",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(domain.Input)),
			fx.ResultTags(`group:"inputs"`),
		),
	),
)

type Web struct {
}

func New(config config.Config) *Web {
	return &Web{}
}

func (s *Web) Name() string {
	return "inputs.web"
}

func (s *Web) Copy(cfg map[string]any) domain.Input {
	return &Web{}
}

func (stages *Web) Do(ctx context.Context, state *domain.State) bool {
	// нужно получить данные пользователя из запроса

	state.User = &domain.User{}
	state.Device = &domain.Device{}

	return true
}
