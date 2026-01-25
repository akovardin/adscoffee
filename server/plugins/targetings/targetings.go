package targetings

import "go.ads.coffee/server/domain"

type Targetings struct {
	list []domain.Targeting
}

func New(list []domain.Targeting) *Targetings {
	return &Targetings{
		list: list,
	}
}

func (i *Targetings) Get(name string, cfg map[string]any) domain.Targeting {
	return nil
}
