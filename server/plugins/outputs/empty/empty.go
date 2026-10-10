package empty

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"outputs.empty",

	fx.Provide(
		fx.Annotate(
			New,
			fx.ResultTags(`group:"outputs"`),
		),
	),
)

type Empty struct {
}

func New() plugins.Provider[plugins.Output] {
	return plugins.Provider[plugins.Output]{
		Name: "outputs.empty",
		Factory: func(cfg map[string]any) (plugins.Output, error) {
			return &Empty{}, nil
		},
	}
}

func (r *Empty) Name() string {
	return "outputs.empty"
}

func (rtb *Empty) Do(ctx context.Context, state *plugins.State) error {
	return nil
}
