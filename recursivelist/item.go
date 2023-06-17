package recursivelist

import (
	// "fmt"

	"github.com/charmbracelet/bubbles/list"
)

type ListItem[T Indentable[U], U ListItem] struct {
	expanded    bool
	Value       T
	CurrentP *ListItem[T]
	ParentModel *Model[T, U]
}

type Indentable[T list.Item] interface {
	list.Item
	Lvl() int
	IndexWithinParent() int
	FilterValue() string
	Children() []Indentable[T]
	Parent() *T
	Add(T)
	AddMulti(...T)
	Find(T) int
	Flatten() []T
	RModify(func(T))
	// ParentOptions() Options
}

func sliceNFind[T Indentable[T]](cur []Indentable[T]) int {
	accum := 0
	for _, i := range cur {
		accum += 1
		accum += sliceNFind[T](i.Children())
	}
	return accum
}

func (i *ListItem[T]) Add(item T, index int) {
	listItem := ListItem[T]{
		Value:       item,
		ParentModel: i.ParentModel,
	}
	listItem.CurrentP = &listItem
	// i.Children() = append(i.Children(), listItem)

	var top *T = item.Parent()
	accum := item.IndexWithinParent()
	for top != nil {
		t := *top
		iwi := t.IndexWithinParent()
		accum += iwi + 1
		if iwi > 0 {
			for _, wee := range (*t.Parent()).Children()[0:iwi - 1] {
				accum += sliceNFind[T](wee.Children())
			}
		}
		top = t.Parent()

	}
	(*item.Parent()).Children()
	i.ParentModel.list.InsertItem(accum + index, listItem)
}

func (i ListItem[T]) Flatten() []T {
	return i.Value.Flatten()
}

func (i ListItem[T]) FilterValue() string {
	return i.Value.FilterValue()
}
func (i ListItem[T]) Expanded() bool {
	return i.expanded
}

func (i ListItem[T]) Point() *ListItem[T] {
	return &i
}

func (i *ListItem[T]) SetExpanded(v bool, pm Model[T]) {
	if pm.Options.Expandable {
		i.expanded = v
	}
}