package gormrelay

type Option[T any] func(*Options[T])

type Options[T any] struct {
	Computed *Computed[T]
}
