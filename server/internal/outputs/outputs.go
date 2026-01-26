package outputs

import "go.ads.coffee/platform/server/domain"

type Outputs struct {
	list map[string]domain.Output
}

func New(list []domain.Output) *Outputs {
	plugins := map[string]domain.Output{}
	for _, output := range list {
		plugins[output.Name()] = output
	}

	return &Outputs{
		list: plugins,
	}
}

func (i *Outputs) Get(name string, cfg map[string]any) domain.Output {
	return i.list[name].Copy(cfg)
}
