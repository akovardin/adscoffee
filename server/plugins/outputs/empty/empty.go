package empty

import (
	"context"

	"go.ads.coffee/platform/server/domain"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"outputs.empty",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(domain.Output)),
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

func (r *Empty) Copy(cfg map[string]any) domain.Output {
	return &Empty{}
}

func (r *Empty) Formats(ff []domain.Format) {
}

func (rtb *Empty) Do(ctx context.Context, state *domain.State) {

	// возвращаем пиксель
}
