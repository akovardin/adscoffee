package analytics

import "context"

type Analytics struct {
}

func New() *Analytics {
	return &Analytics{}
}

func (a *Analytics) LogImpression(ctx context.Context) {

}

func (a *Analytics) LogClick(ctx context.Context) {

}
