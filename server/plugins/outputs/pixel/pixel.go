package pixel

import (
	"context"
	"image"
	"image/color"
	"image/gif"
	"net/http"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"outputs.pixel",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Output)),
			fx.ResultTags(`group:"outputs"`),
		),
	),
)

type Pixel struct {
}

func New() *Pixel {
	return &Pixel{}
}

func (r *Pixel) Name() string {
	return "outputs.pixel"
}

func (r *Pixel) Copy(cfg map[string]any) plugins.Output {
	return &Pixel{}
}

func (rtb *Pixel) Do(ctx context.Context, state *plugins.State) error {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{255, 0, 0, 255})
	if err := gif.Encode(state.Response, img, nil); err != nil {

		state.Response.WriteHeader(http.StatusInternalServerError)

		return err
	}

	state.Response.Header().Add("Content-Type", "image/gif")
	state.Response.WriteHeader(http.StatusOK)

	return nil
}
