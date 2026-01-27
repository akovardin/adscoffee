package rotation

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"stages.rotation",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Stage)),
			fx.ResultTags(`group:"stages"`),
		),
	),
)

type Rotattion struct{}

func New() *Rotattion {
	return &Rotattion{}
}

func (t *Rotattion) Name() string {
	return "stages.rotation"
}

func (t *Rotattion) Copy(cfg map[string]any) plugins.Stage {
	return &Rotattion{}
}

func (t *Rotattion) Do(ctx context.Context, state *plugins.State) {
	// делаем взвешанный рандом по ecpm
	state.Winners = state.Candidates
}
