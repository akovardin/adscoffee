package formats

import (
	"context"
	"fmt"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

const TypeBanner = "banner"

type Banner struct {
	base string
}

func NewBanner() *Banner {
	return &Banner{}
}

func (b *Banner) Name() string {
	return "banner"
}

func (b *Banner) Copy(cfg map[string]any) plugins.Format {
	base, _ := cfg["base"].(string)

	return &Banner{
		base: base,
	}
}

func (b *Banner) Render(ctx context.Context, state *plugins.State) (any, error) {
	for _, b := range state.Winners {
		if b.Format != TypeBanner {
			continue
		}

		return fmt.Sprintf(`<a href="%s"><img src="%s"></a>`, b.Target, b.Image.Full("example")), nil
	}

	return nil, nil
}
