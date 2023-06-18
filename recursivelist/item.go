package recursivelist

import (
	// "fmt"

	"math"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ListItem[T ItemWrapper[T]] struct {
	expanded    bool
	value       *T
	ParentModel *Model[T]
	Children    *[]ListItem[T]
	Parent      *ListItem[T]
}

type ItemWrapper[T any] interface {
	list.Item
	Find(T) int
	IndexWithinParent() int
	GetChildren() []ItemWrapper[T]
	SetChildren([]ItemWrapper[T])
	GetParent() *T
	Value() *T
	Lvl() int
}

type IIndentable[T ItemWrapper[T]] interface {
	list.Item
	FilterValue() string
	Add(IIndentable[T])
	AddMulti(...T)
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
	choild := make([]ItemWrapper[T], 0)
	for _, val := range k {
		choild = append(choild, *val.value)
	}
	(*r.value).SetChildren(choild)
	*r.Children = k
}

func (r ListItem[T]) FilterValue() string {
	return (*r.value).FilterValue()
}

func (r ListItem[T]) Value() T {
	return *r.value
}

func (r ListItem[T]) Flatten() []ListItem[T] {
	accum := make([]ListItem[T], 0)
	accum = append(accum, r)
	for _, ite := range *r.Children {
		accum = append(accum, ite.Flatten()...)
	}
	return accum
}

func (r ListItem[T]) RModify(fnn func(ListItem[T])) {
	for _, val := range *r.Children {
		// val.RModify(fnn)
		fnn(val)
	}
}

func (r ListItem[T]) GetChildren() []ListItem[T] {
	ret := make([]ListItem[T], 0)
	for _, i := range *r.Children {
		ret = append(ret, i)
	}
	return ret
}

func (r ListItem[T]) GetParent() *ListItem[T] {
	return r.Parent
}

func (r ListItem[T]) TotalBeneath() int {
	if r.Parent != nil {
		v := r.Parent.point()
		return (v.Find(*r.value))
	}

	return 0
}

func (r ListItem[T]) IndexWithinParent() int {
	if r.Parent != nil {
		v := r.Parent.point()
		return (v.Find(*r.value))
	}

	return 0
}

func sliceNFind[T ItemWrapper[T]](cur []ItemWrapper[T]) int {
	accum := 0
	for _, i := range cur {
		accum += 1
		accum += sliceNFind[T](i.GetChildren())
	}
	return accum
}

func (i *ListItem[T]) realAdd(item T, index int) {
	whatthefuck := make([]ListItem[T], 0)
	listItem := ListItem[T]{
		value:       &item,
		ParentModel: i.ParentModel,
		Children:    &whatthefuck,
		Parent:      i,
	}
	
	*i.Children = append(*i.Children, listItem)

	var top *T = listItem.Parent.value
	accum := item.IndexWithinParent()
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
	i.ParentModel.list.InsertItem(accum+index, listItem)
}

func (i ListItem[T]) Add(item T, index int) {
	point := &i
	point.realAdd(item, index)
}

func (i ListItem[T]) Expanded() bool {
	return i.expanded
}

func (i ListItem[T]) Point() *ListItem[T] {
	return &i
}

func (i *ListItem[T]) SetExpanded(v bool) tea.Cmd {
	if i.ParentModel.Options.Expandable {
		i.expanded = v
	}
	if !v {
		for ii := 0; ii < i.TotalBeneath(); ii++ {
			i.ParentModel.list.RemoveItem(ii)
		}
	}
	return nil
}

func NewItem[T ItemWrapper[T]](item T, del list.ItemDelegate) ListItem[T] {
	childVar := make([]ListItem[T], 0)
	li := ListItem[T]{
		value:       &item,
		ParentModel: &Model[T]{},
		expanded:    true,
		Children:    &childVar,
	}
	li.ParentModel.list = list.New([]list.Item{}, del, 200, 200)
	return li
}
