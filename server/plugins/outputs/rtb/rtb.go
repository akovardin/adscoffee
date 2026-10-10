package rtb

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"outputs.rtb",

	fx.Provide(
		fx.Annotate(
			New,
			fx.ResultTags(`group:"outputs"`),
		),
	),
)

type Rtb struct {
}

func New() plugins.Provider[plugins.Output] {
	return plugins.Provider[plugins.Output]{
		Name: "outputs.rtb",
		Factory: func(cfg map[string]any) (plugins.Output, error) {
			return &Rtb{}, nil
		},
	}
}

func (r *Rtb) Name() string {
	return "outputs.rtb"
}

func (rtb *Rtb) Do(ctx context.Context, state *plugins.State) error {
	return nil
}
