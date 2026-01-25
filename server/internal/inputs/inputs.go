package inputs

import (
	"go.ads.coffee/server/domain"
)

type Inputs struct {
	plugins map[string]domain.Input
}

func New(inputs []domain.Input) *Inputs {
	plugins := map[string]domain.Input{}
	for _, input := range inputs {
		plugins[input.Name()] = input
	}

	return &Inputs{
		plugins: plugins,
	}
}

func (i *Inputs) Get(name string, cfg map[string]any) domain.Input {
	return i.plugins[name].Copy(cfg)
}
