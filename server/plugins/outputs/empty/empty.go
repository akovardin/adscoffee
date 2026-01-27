package empty

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"outputs.empty",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Output)),
			fx.ResultTags(`group:"outputs"`),
		),
	),
)

type Empty struct {
}

func New() *Empty {
	return &Empty{}
}

func (r *Empty) Name() string {
	return "outputs.empty"
}

func (r *Empty) Copy(cfg map[string]any) plugins.Output {
	return &Empty{}
}

func (r *Empty) Formats(ff []plugins.Format) {
}

func (rtb *Empty) Do(ctx context.Context, state *plugins.State) {

	// возвращаем пиксель
}
