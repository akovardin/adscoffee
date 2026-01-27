package banner

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/ads"
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

func (b *Banner) Render(ctx context.Context, state *plugins.State) error {
	items := []ads.Banner{}

	for _, b := range state.Winners {
		if b.Format != "banner" {
			continue
		}

		items = append(items, b)
	}

	data, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("error on marshal banner: %w", err)
	}

	state.Response.Write(data)

	return nil
}
