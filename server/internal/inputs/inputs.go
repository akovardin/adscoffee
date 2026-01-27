package inputs

import (
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

type Inputs struct {
	plugins map[string]plugins.Input
}

func New(inputs []plugins.Input) *Inputs {
	plugins := map[string]plugins.Input{}
	for _, input := range inputs {
		plugins[input.Name()] = input
	}

	return &Inputs{
		plugins: plugins,
	}
}

func (i *Inputs) Get(name string, cfg map[string]any) plugins.Input {
	return i.plugins[name].Copy(cfg)
}
