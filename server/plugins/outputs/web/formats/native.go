package formats

import (
	"context"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

const TypeNative = "native"

type Native struct{}

func NewNative() *Native {
	return &Native{}
}

func (b *Native) Name() string {
	return "native"
}

type NativeResponse struct {
	Title string `json:"title"`
	Img   string `json:"img"`
}

func (b *Native) Render(ctx context.Context, state *plugins.State) (any, error) {
	items := []NativeResponse{}

	for _, b := range state.Winners {
		if b.Format != TypeNative {
			continue
		}

		items = append(items, NativeResponse{
			Title: b.Title,
			Img:   b.Image.Full("example"),
		})
	}

	return items, nil
}
