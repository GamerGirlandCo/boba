package utils

import (
	"reflect"
	"time"
)

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

func FindField(val interface{}, name string) (*reflect.Value, int) {
	var retI int = -1
	fields := reflect.TypeOf(val)
	values := reflect.ValueOf(val)
	for i := 0; i < fields.NumField(); i++ {
		cf := fields.Field(i).Name
		if cf == name {
			fv := values.Field(i)
			return &fv, i
		}
	}
	fv := reflect.ValueOf(nil)
	return &fv, retI
}



func FieldByName(val interface{}, name string) (reflect.Value, reflect.Type) {
	// var retI int = -1
	// fields := reflect.TypeOf(val)
	values := reflect.ValueOf(val)
	fav := reflect.Indirect(values).FieldByName(name)
	// fff := reflect.TypeOf(fav)
	return fav, fav.Type()
}

func EqualIndex(a []interface{}, b interface{}) int {
	for i := range a {
		if reflect.DeepEqual(a[i], b) {
			return i
		}
	}
	return -1
}