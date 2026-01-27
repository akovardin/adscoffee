package targetings

import "go.ads.coffee/platform/server/internal/domain/plugins"

type Targetings struct {
	list map[string]plugins.Targeting
}

func New(list []plugins.Targeting) *Targetings {
	plugins := map[string]plugins.Targeting{}
	for _, targeting := range list {
		plugins[targeting.Name()] = targeting
	}

	return &Targetings{
		list: plugins,
	}
}

func (i *Targetings) Get(name string, cfg map[string]any) plugins.Targeting {
	return i.list[name].Copy(cfg)
}
