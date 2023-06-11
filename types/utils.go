package types

import "time"

type Predicate[T any] func(t T) bool
type GenResultMsg[T any] struct {
	Res T
}
type TickMsg time.Time