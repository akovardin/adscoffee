package stages

import "go.ads.coffee/platform/server/internal/domain/plugins"

type Stages struct {
	list map[string]plugins.Stage
}

func New(list []plugins.Stage) *Stages {
	plugins := map[string]plugins.Stage{}
	for _, stage := range list {
		plugins[stage.Name()] = stage
	}

	return &Stages{
		list: plugins,
	}
}

func (i *Stages) Get(name string, cfg map[string]any) plugins.Stage {
	return i.list[name].Copy(cfg)
}
