package banner

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"formats.banner",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Format)),
			fx.ResultTags(`group:"formats"`),
		),
	),
)

type Banner struct{}

func New() *Banner {
	return &Banner{}
}

func (b *Banner) Name() string {
	return "formats.banner"
}

func (b *Banner) Copy(cfg map[string]any) plugins.Format {
	return &Banner{}
}

func (b *Banner) Render(ctx context.Context, state *plugins.State) {

}
