package pixel

import (
	"context"

	"go.ads.coffee/platform/server/domain"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"outputs.pixel",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(domain.Output)),
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

func (r *Pixel) Copy(cfg map[string]any) domain.Output {
	return &Pixel{}
}

func (r *Pixel) Formats(ff []domain.Format) {
}

func (rtb *Pixel) Do(ctx context.Context, state *domain.State) {

	// возвращаем пиксель
}
