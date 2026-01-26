package targetings

import "go.ads.coffee/platform/server/domain"

type Targetings struct {
	list map[string]domain.Targeting
}

func New(list []domain.Targeting) *Targetings {
	plugins := map[string]domain.Targeting{}
	for _, targeting := range list {
		plugins[targeting.Name()] = targeting
	}

	return &Targetings{
		list: plugins,
	}
}

func (i *Targetings) Get(name string, cfg map[string]any) domain.Targeting {
	return i.list[name].Copy(cfg)
}
