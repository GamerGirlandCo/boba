package recursivelist

import (
	"git.tablet.sh/tablet/boba/utils"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ListItem[T ItemWrapper[T]] struct {
	// whether or not this item is expanded.
	expanded *bool
	// the raw value that this item points to
	value       *T
	ParentModel *Model[T]
	// this is a recursive struct, meaning that it can
	// have 0 or more child elements, which can have 0
	// or more child elements, and so on
	Children *[]ListItem[T]
	// if it is a top level item, this field will be `nil`.
	// otherwise contains the element that contains the
	// current item
	Parent *ListItem[T]
}

type ItemWrapper[T any] interface {
	list.Item
	// finds and returns the index of the child
	// within the current element's children
	Find(T) int
	// should find and return the index of the
	// current element within its parent
	IndexWithinParent() int
	// should return this item's children
	GetChildren() []T
	// overwrites the value's children with the given argument.
	SetChildren([]T)
	// should return the current element's parent.
	GetParent() *T
	// should set the item's parent to the argument.
	SetParent(*T)
	// since `value` is not exported from the `ListItem`,
	// we need a wrapper function to access it.
	Value() *T
	// returns how deeply nested the current node is, as an int.
	// i'd recommend calculating this by checking that `Parent`
	// is not null in a for loop.
	Lvl() int
	// adds an element to this item's children.
	Add(int, T)
	// same as Add() but can receive more than one argument!
	AddMulti(int, ...T)
	// Removes child at specified index and returns it
	Remove(int) T
	// Returns a string representation of this item.
	ToString() string
}

func (r ListItem[T]) point() T {
	return *r.value
}

func (r ListItem[T]) FilterValue() string {
	return (*r.value).FilterValue()
}

func (r *ListItem[T]) Value() *T {
	return r.value
}

func (r ListItem[T]) Flatten() []ListItem[T] {
	accum := make([]ListItem[T], 0)
	accum = append(accum, r)
	for _, ite := range *r.Children {
		if *ite.expanded {
			accum = append(accum, ite.Flatten()...)
		} else {
			accum = append(accum, ite)
		}
	}
	return accum
}

func (r ListItem[T]) RModify(fnn func(ListItem[T])) {
	fnn(r)
	for _, val := range *r.Children {
		val.RModify(fnn)
	}
}

func (r ListItem[T]) GetChildren() []ListItem[T] {
	ret := make([]ListItem[T], 0)
	ret = append(ret, *r.Children...)
	return ret
}

func (r ListItem[T]) GetParent() *ListItem[T] {
	return r.Parent
}

func (r *ListItem[T]) setParent(pa *T) {
	(*pa).SetParent(r.Parent.value)
	*r.Parent = r.ParentModel.NewItem(*pa)
}

func (r ListItem[T]) TotalBeneath() int {
	accum := 0
	for _, val := range *r.Children {
		accum += 1
		accum += val.TotalBeneath()
	}
	return accum
}

func (r ListItem[T]) IndexWithinParent() int {
	// you're supposed to implement this yourself...
	if r.value != nil {
		return (*r.value).IndexWithinParent()
	}
	return 0
}

func (r ListItem[T]) everythingBefore() int {
	a := 1
	top := r.Parent
	for top != nil {
		a++
		a += top.IndexWithinParent()
		top = (*top).Parent
	}
	return a
}

func sliceNFind[T ItemWrapper[T]](cur []T) int {
	accum := 0
	for _, i := range cur {
		accum += 1
		accum += sliceNFind[T]((*(*i.Value()).Value()).GetChildren())
	}
	return accum
}

// since go will complain about `Add()` having a pointer receiver...
func (i *ListItem[T]) realAdd(arg T, index int) ListItem[T] {

	item := i.ParentModel.NewItem(arg)
	var nindex int
	if index >= len(*i.Children) {
		*i.Children = append(*i.Children, item)
	} else {
		nindex = utils.MaxInt(0, index)
		*i.Children = utils.SliceInsert[ListItem[T]](*i.Children, nindex, item)
	}
	(*i.value).Add(nindex, arg)

	item.ParentModel = i.ParentModel
	_, accu := item.ParentModel.Flatten()
	item.ParentModel.List.SetItems(accu)
	return item
}

func (i *ListItem[T]) Add(item T, index int) ListItem[T] {
	//	(*i.value).Add(index, item)
	return i.realAdd(item, index)
}

func (i *ListItem[T]) AddMulti(index int, items ...T) {
	//(*i.value).AddMulti(index, items...)
	for mi, val := range items {
		i.Add(val, index+mi+1)
	}
}

func (i ListItem[T]) Expanded() bool {
	return *i.expanded
}

func (i ListItem[T]) Point() *ListItem[T] {
	return &i
}

func (i *ListItem[T]) SetExpanded(v bool) tea.Cmd {
	if i.ParentModel.Options.Expandable {
		*i.expanded = v
	}
	return tea.EnterAltScreen
}
