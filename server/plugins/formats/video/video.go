package video

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"formats.video",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Format)),
			fx.ResultTags(`group:"formats"`),
		),
	),
)

// Example plugin
type Video struct{}

func New() *Video {
	return &Video{}
}

func (b *Video) Name() string {
	return "formats.video"
}

func (b *Video) Copy(cfg map[string]any) plugins.Format {
	return &Video{}
}

func (b *Video) Render(ctx context.Context, state *plugins.State) error {
	return nil
}
