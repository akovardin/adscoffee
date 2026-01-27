package native

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"formtas.native",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Format)),
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

func (b *Native) Copy(cfg map[string]any) plugins.Format {
	return &Native{}
}

func (b *Native) Render(ctx context.Context, state *plugins.State) {

}
