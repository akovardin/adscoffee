package rotation

import (
	"context"

	"go.ads.coffee/platform/server/domain"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"stages.rotation",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(domain.Stage)),
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

func (t *Rotattion) Copy(cfg map[string]any) domain.Stage {
	return &Rotattion{}
}

func (t *Rotattion) Do(ctx context.Context, state *domain.State) {
	// делаем взвешанный рандом по ecpm
	state.Winners = state.Candidates
}
