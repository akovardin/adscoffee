package static

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"inputs.static",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Input)),
			fx.ResultTags(`group:"inputs"`),
		),
	),
)

type Static struct {
}

func New() *Static {
	return &Static{}
}

func (s *Static) Name() string {
	return "inputs.static"
}

func (s *Static) Copy(cfg map[string]any) plugins.Input {
	return &Static{}
}

func (stages *Static) Do(ctx context.Context, state *plugins.State) bool {

	return true
}
