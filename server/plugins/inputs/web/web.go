package web

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"inputs.web",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Input)),
			fx.ResultTags(`group:"inputs"`),
		),
	),
)

type Web struct {
}

func New() *Web {
	return &Web{}
}

func (s *Web) Name() string {
	return "inputs.web"
}

func (s *Web) Copy(cfg map[string]any) plugins.Input {
	return &Web{}
}

func (stages *Web) Do(ctx context.Context, state *plugins.State) bool {
	// нужно получить данные пользователя из запроса

	state.User = &plugins.User{}
	state.Device = &plugins.Device{}

	// проверить наличие юнита

	// по юниту может быть включена медиация ->
	// это значит, что должны получить разные настройки сети и выполнить
	// аукцион и ротацию баннера

	return true
}
