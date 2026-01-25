package domain

type Stage interface {
	Name() string
}

type WithTargetings interface {
	Targetings(tt []Targeting)
}
