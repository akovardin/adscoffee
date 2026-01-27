package web

import (
	"context"
	"fmt"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"outputs.web",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Output)),
			fx.ResultTags(`group:"outputs"`),
		),
	),
)

type Web struct {
	formtats []plugins.Format
}

func New() *Web {
	return &Web{}
}

func (w *Web) Name() string {
	return "outputs.web"
}

func (w *Web) Copy(cfg map[string]any) plugins.Output {
	return &Web{
		formtats: w.formtats,
	}
}

func (w *Web) Formats(ff []plugins.Format) {
	w.formtats = ff
}

//nolint:errcheck
func (w *Web) Do(ctx context.Context, state *plugins.State) error {
	for _, f := range w.formtats {
		if err := f.Render(ctx, state); err != nil {
			return fmt.Errorf("error on render format: %w", err)
		}
	}

	return nil
}
