package pixel

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"outputs.pixel",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Output)),
			fx.ResultTags(`group:"outputs"`),
		),
	),
)

type Pixel struct {
}

func New() *Pixel {
	return &Pixel{}
}

func (r *Pixel) Name() string {
	return "outputs.pixel"
}

func (r *Pixel) Copy(cfg map[string]any) plugins.Output {
	return &Pixel{}
}

func (rtb *Pixel) Do(ctx context.Context, state *plugins.State) error {
	return nil
}
