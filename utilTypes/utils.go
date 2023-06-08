package utilTypes
type Predicate[T any] func(t T) bool
type GenResultMsg struct {
	res string
}