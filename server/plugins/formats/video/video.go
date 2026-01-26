package video

import (
	"context"

	"go.ads.coffee/platform/server/domain"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"formats.video",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(domain.Format)),
			fx.ResultTags(`group:"formats"`),
		),
	),
)

type Video struct{}

func New() *Video {
	return &Video{}
}

func (b *Video) Name() string {
	return "formats.video"
}

func (b *Video) Copy(cfg map[string]any) domain.Format {
	return &Video{}
}

func (b *Video) Render(ctx context.Context, state *domain.State) {

}
