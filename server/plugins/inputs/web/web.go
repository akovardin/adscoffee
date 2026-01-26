package web

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/server/config"
	"go.ads.coffee/server/domain"
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
	input := &Web{}

	return input
}

func (s *Web) Name() string {
	return "inputs.web"
}

func (s *Web) Copy(cfg map[string]any) domain.Input {
	return &Web{} // copy
}

func (stages *Web) Do(ctx context.Context, state *domain.State) bool {

	// обработка разных типов запросов тоже
	// может быть вынесена в пллагины

	return true
}
