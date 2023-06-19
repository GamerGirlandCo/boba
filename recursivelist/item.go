package recursivelist

import (
	"math"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ListItem[T ItemWrapper[T]] struct {
	// whether or not this item is expanded.
	expanded    *bool
	// the raw value that this item points to
	value       *T
	ParentModel *Model[T]
	// this is a recursive struct, meaning that it can
	// have 0 or more child elements, which can have 0
	// or more child elements, and so on
	Children    []ListItem[T]
	// if it is a top level item, this field will be `nil`.
	Parent      *ListItem[T]
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
	// since `value` is not exported from the `ListItem`,
	// we need a wrapper function to access it.
	Value() *T
	// returns how deeply nested the current node is, as an int.
	// i'd recommend calculating this by checking that `Parent` 
	// is not null in a for loop.
	Lvl() int
	// adds an element to this item's children.
	Add(T)
	// same as Add() but can receive more than one argument!
	AddMulti(...T)
}

type IIndentable[T ItemWrapper[T]] interface {
	list.Item
	FilterValue() string
	Flatten() []ListItem[T]
	RModify(func(T))
	Value() *T
	ParentOptions() Options
	SetOptions(Options)
	TotalBeneath() int
}

func (r ListItem[T]) point() T {
	return *r.value
}

func (r *ListItem[T]) SetChildren(k []ListItem[T]) {
	choild := make([]T, 0)
	for _, val := range k {
		choild = append(choild, *val.value)
	}
	(*r.value).SetChildren(choild)
	r.Children = k
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
	for _, ite := range r.Children {
		if *ite.expanded {
			accum = append(accum, ite.Flatten()...)
		}  else {
			accum = append(accum, ite)
		}
	}
	return accum
}

func (r ListItem[T]) RModify(fnn func(ListItem[T])) {
	for _, val := range r.Children {
		// val.RModify(fnn)
		fnn(val)
	}
}

func (r ListItem[T]) GetChildren() []ListItem[T] {
	ret := make([]ListItem[T], 0)
	ret = append(ret, r.Children...)
	return ret
}

func (r ListItem[T]) GetParent() *ListItem[T] {
	return r.Parent
}

func (r ListItem[T]) TotalBeneath() int {
	accum := 0
	for _, val := range r.Children {
		accum += 1
		accum += val.TotalBeneath()
	}
	return accum
}

func (r ListItem[T]) IndexWithinParent() int {
	// you're supposed to implement this yourself...
	return (*r.value).IndexWithinParent()
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
func (i *ListItem[T]) realAdd(item ListItem[T], index int) {
	
	i.Children = append(i.Children, item)
	var top *T
	if item.Parent != nil {
		top = item.Parent.value
	}
	item.ParentModel = i.ParentModel
	accum := item.everythingBefore()
	for top != nil {
		iwi := (*top).IndexWithinParent()
		if iwi >= 0 {
			accum += iwi + 1
			slic := int(math.Min(
				float64(iwi),
				float64(len((*top).GetChildren())),
			))
			p := (*top)
			for _, wee := range (p).GetChildren()[0:slic] {
				accum += sliceNFind[T](wee.GetChildren())
			}
		}

		top = (*top).GetParent()
	}
	i.ParentModel.List.InsertItem(accum+index, item)
}

func (i ListItem[T]) Add(item ListItem[T], index int) {
	(&i).realAdd(item, index)
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

// creates a new ListItem
func NewItem[T ItemWrapper[T]](item T, pm Model[T]) ListItem[T] {
	childVar := make([]ListItem[T], 0)
	expanded := true
	li := ListItem[T]{
		value: &item,
		ParentModel: &Model[T]{},
		Children: childVar,
		expanded: &expanded,
	}
	*li.ParentModel = pm
	return li
}
