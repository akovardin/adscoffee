package rtb

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"outputs.rtb",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Output)),
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

func (r *Rtb) Copy(cfg map[string]any) plugins.Output {
	return &Rtb{}
}

func (r *Rtb) Formats(ff []plugins.Format) {
	// set formats
}

func (rtb *Rtb) Do(ctx context.Context, state *plugins.State) {

	// обработка разных типов запросов тоже
	// может быть вынесена в пллагины
}
