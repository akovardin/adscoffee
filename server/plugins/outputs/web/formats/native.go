package formats

import (
	"context"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

const TypeNative = "native"

type Native struct {
	base string
}

func NewNative() *Native {
	return &Native{}
}

func (b *Native) Name() string {
	return "native"
}

type NativeResponse struct {
	Description string   `json:"description"`
	Title       string   `json:"information"`
	Image       string   `json:"image"`
	Target      string   `json:"target"`
	Impressions []string `json:"impressions"`
	Clicks      []string `json:"click"`
}

func (b *Native) Copy(cfg map[string]any) plugins.Format {
	base, _ := cfg["base"].(string)

	return &Native{
		base: base,
	}
}

func (b *Native) Render(ctx context.Context, state *plugins.State) (any, error) {
	items := []NativeResponse{}

	for _, b := range state.Winners {
		items = append(items, NativeResponse{
			Title:       b.Title,
			Description: b.Description,
			Target:      b.Target,
			Image:       b.Image.Full(""),

			Impressions: []string{}, // TODO: make trackers
			Clicks:      []string{},
		})
	}

	return items, nil
}
