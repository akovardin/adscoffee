package stages

import "go.ads.coffee/server/domain"

type Stages struct {
	list map[string]domain.Stage
}

func New(list []domain.Stage) *Stages {
	plugins := map[string]domain.Stage{}
	for _, stage := range list {
		plugins[stage.Name()] = stage
	}

	return &Stages{
		list: plugins,
	}
}

func (i *Stages) Get(name string, cfg map[string]any) domain.Stage {
	return i.list[name].Copy(cfg)
}
