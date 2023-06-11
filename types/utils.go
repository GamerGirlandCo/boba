package types

import "time"

type Predicate[T any] func(t T) bool

type IToString interface {
	ToString() string
}

type Capsule[T any] struct {
	value T
}

type GenResultMsg[T any] struct {
	Res T
}

type TickMsg time.Time
