package outputs

import "go.ads.coffee/server/domain"

type Outputs struct {
	list []domain.Output
}

func New(list []domain.Output) *Outputs {
	return &Outputs{
		list: list,
	}
}

func (i *Outputs) Get(name string, cfg map[string]any) domain.Output {
	return nil
}
