package rtb

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/plugins"
	"go.ads.coffee/platform/server/internal/metrics"
)

var Module = fx.Module(
	"inputs.rtb",

	fx.Provide(
		fx.Annotate(
			New,
			fx.ResultTags(`group:"inputs"`),
		),
	),
)

type Rtb struct {
}

func New() plugins.Provider[plugins.Input] {
	return plugins.Provider[plugins.Input]{
		Name: "inputs.rtb",
		Factory: func(cfg map[string]any) (plugins.Input, error) {
			return &Rtb{}, nil
		},
	}
}

func (rtb *Rtb) Name() string {
	return "inputs.rtb"
}

func (rtb *Rtb) Do(ctx context.Context, state *plugins.State) bool {
	metrics.IncInput(rtb.Name())

	// обработка разных типов запросов тоже
	// может быть вынесена в пллагины

	// тут я могу понять какие форматы мне нужны

	// разбираем rtb запрос

	return true
}
