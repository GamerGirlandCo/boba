package utils

import (
	"reflect"
	"time"

	"github.com/charmbracelet/bubbles/key"
)

type Predicate[T any] func(t T) bool

type Stringable interface {
	ToString() string
}

type Capsule[T any] struct {
	value T
}

type GenResultMsg struct {
	Res any
	StringRep string
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

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func IterKeybindings(v ...interface{}) []key.Binding {
	kbt := reflect.ValueOf(key.NewBinding())
	var rv []key.Binding
	for _, i := range v {
		// f := reflect.TypeOf(i)
		voi := reflect.ValueOf(i)
		for j := 0; j < voi.NumField(); j++ {
			if voi.Field(j).Type().Name() == kbt.Type().Name() {
				rv = append(rv, voi.Field(j).Interface().(key.Binding))
			}
		}
	}
	return rv
}