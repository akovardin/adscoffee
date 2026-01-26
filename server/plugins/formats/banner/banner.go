package banner

import (
	"context"

	"go.ads.coffee/platform/server/domain"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"formats.banner",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(domain.Format)),
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

func (b *Banner) Copy(cfg map[string]any) domain.Format {
	return &Banner{}
}

func (b *Banner) Render(ctx context.Context, state *domain.State) {

}
