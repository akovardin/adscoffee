package formats

import (
	"context"

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

type BannerResponse struct {
	Title string `json:"title"`
	Img   string `json:"img"`
}

func (b *Banner) Copy(cfg map[string]any) plugins.Format {
	base, _ := cfg["base"].(string)

	return &Banner{
		base: base,
	}
}

func (b *Banner) Render(ctx context.Context, state *plugins.State) (any, error) {
	items := []BannerResponse{}

	for _, b := range state.Winners {
		items = append(items, BannerResponse{
			Title: b.Title,
			Img:   b.Image.Full("example"),
		})
	}

	return items, nil
}
