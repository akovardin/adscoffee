package rtb

import (
	"context"

	"go.ads.coffee/server/domain"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"outputs.rtb",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(domain.Output)),
			fx.ResultTags(`group:"outputs"`),
		),
	),
)

type Rtb struct {
}

func New() *Rtb {
	return &Rtb{}
}

func (r *Rtb) Name() string {
	return "outputs.rtb"
}

func (r *Rtb) Copy(cfg map[string]any) domain.Output {
	return &Rtb{}
}

func (r *Rtb) Formats(ff []domain.Format) {
	// set formats
}

func (rtb *Rtb) Do(ctx context.Context, state *domain.State) {

	// обработка разных типов запросов тоже
	// может быть вынесена в пллагины
}
