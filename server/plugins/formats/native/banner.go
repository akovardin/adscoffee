package native

import (
	"context"

	"go.ads.coffee/server/domain"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"formtas.native",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(domain.Format)),
			fx.ResultTags(`group:"formats"`),
		),
	),
)

type Native struct{}

func New() *Native {
	return &Native{}
}

func (b *Native) Name() string {
	return "formats.native"
}

func (b *Native) Copy(cfg map[string]any) domain.Format {
	return &Native{}
}

func (b *Native) Render(ctx context.Context, state *domain.State) {

}
