package stages

import "go.ads.coffee/server/domain"

type Stages struct {
	list []domain.Stage
}

func New(list []domain.Stage) *Stages {
	return &Stages{
		list: list,
	}
}

func (i *Stages) Get(name string, cfg map[string]any) domain.Stage {
	return nil
}
