package targeting

import "go.ads.coffee/server/domain"

type Targeting struct{}

func New() *Targeting {
	return &Targeting{}
}

func (t *Targeting) Targetings(tt []domain.Targeting) {

}

func (t *Targeting) Name() string {
	return "targetings"
}
