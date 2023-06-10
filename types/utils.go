package types

type Predicate[T any] func(t T) bool
type GenResultMsg[T any] struct {
	Res T
}
