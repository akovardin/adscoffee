package plugins

import "context"

type Stage interface {
	Name() string
	Do(ctx context.Context, state *State) error
}

type WithTargetings interface {
	Targetings(tt []Targeting)
}

type WithLimits interface {
	Limits(ll []Limit)
}
