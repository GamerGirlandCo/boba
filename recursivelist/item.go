package recursivelist

import (
	// "fmt"

	// "log"
	"math"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ListItem[T ItemWrapper[T]] struct {
	// whether or not this item is expanded.
	// note that you have to implement your
	// own state management
	// for expansion/collapsing via your delegate's
	// `Update` function
	expanded    *bool
	value       *T
	ParentModel *Model[T]
	Children    []ListItem[T]
	Parent      *ListItem[T]
}

type ItemWrapper[T any] interface {
	list.Item
	Find(T) int
	IndexWithinParent() int
	GetChildren() []T
	SetChildren([]T)
	GetParent() *T
	Value() *T
	Lvl() int
	Add(T)
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
		// accum = append(accum, ite)
		accum = append(accum, ite.Flatten()...)
		// if *ite.expanded {
		// } else {
		// 	accum = append(accum, ite)
		// }
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
	for _, i := range r.Children {
		ret = append(ret, i)
	}
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
	// if r.Parent != nil {
	// 	v := r.Parent.point()
	// 	return (v.Find(*r.value))
	// }
	return (*r.value).IndexWithinParent()
	// return 0
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

func NewItem[T ItemWrapper[T]](item T, del list.ItemDelegate, opts Options, pm Model[T]) ListItem[T] {
	childVar := make([]ListItem[T], 0)
	var thing Options = opts
	if thing.ClosedPrefix == "" || thing.OpenPrefix == "" || thing.Width == 0 || thing.Height == 0 {
		thing = DefaultOptions
	}
	expanded := true
	li := ListItem[T]{
		value: &item,
		ParentModel: &Model[T]{
			Options: thing,
		},
		Children: childVar,
		expanded: &expanded,
	}
	*li.ParentModel = pm
	return li
}
