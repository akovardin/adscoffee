package outputs

import "go.ads.coffee/platform/server/internal/domain/plugins"

type Outputs struct {
	list map[string]plugins.Output
}

func New(list []plugins.Output) *Outputs {
	plugins := map[string]plugins.Output{}
	for _, output := range list {
		plugins[output.Name()] = output
	}

	return &Outputs{
		list: plugins,
	}
}

func (i *Outputs) Get(name string, cfg map[string]any) plugins.Output {
	return i.list[name].Copy(cfg)
}
