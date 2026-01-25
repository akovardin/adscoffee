package domain

type Targeting interface {
	Filter()
	Name() string
}
