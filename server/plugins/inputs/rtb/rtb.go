package rtb

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"inputs.rtb",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Input)),
			fx.ResultTags(`group:"inputs"`),
		),
	),
)

type Rtb struct {
}

func New() *Rtb {
	return &Rtb{}
}

func (rtb *Rtb) Name() string {
	return "inputs.rtb"
}

func (rtb *Rtb) Copy(cfg map[string]any) plugins.Input {
	return &Rtb{}
}

func (rtb *Rtb) Do(ctx context.Context, state *plugins.State) bool {

	// обработка разных типов запросов тоже
	// может быть вынесена в пллагины

	// тут я могу понять какие форматы мне нужны

	// разбираем rtb запрос

	return true
}
