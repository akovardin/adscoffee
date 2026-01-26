package rtb

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/server/config"
	"go.ads.coffee/server/domain"
)

var Module = fx.Module(
	"inputs.rtb",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(domain.Input)),
			fx.ResultTags(`group:"inputs"`),
		),
	),
)

type Rtb struct {
}

func New(config config.Config) *Rtb {
	input := &Rtb{}

	return input
}

func (rtb *Rtb) Name() string {
	return "inputs.rtb"
}

func (rtb *Rtb) Copy(cfg map[string]any) domain.Input {
	return &Rtb{} // copy
}

func (rtb *Rtb) Do(ctx context.Context, state *domain.State) bool {

	// обработка разных типов запросов тоже
	// может быть вынесена в пллагины

	// тут я могу понять какие форматы мне нужны

	return true
}
