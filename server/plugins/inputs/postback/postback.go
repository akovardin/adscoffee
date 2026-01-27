package postback

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"inputs.postback",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Input)),
			fx.ResultTags(`group:"inputs"`),
		),
	),
)

type Postback struct {
}

func New() *Postback {
	return &Postback{}
}

func (s *Postback) Name() string {
	return "inputs.postback"
}

func (s *Postback) Copy(cfg map[string]any) plugins.Input {
	return &Postback{}
}

func (stages *Postback) Do(ctx context.Context, state *plugins.State) bool {
	// нужно получить данные пользователя из запроса

	state.User = &plugins.User{}
	state.Device = &plugins.Device{}

	return true
}
