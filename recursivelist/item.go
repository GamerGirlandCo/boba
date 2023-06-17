package recursivelist

import (
	// "fmt"

	"github.com/charmbracelet/bubbles/list"
)

type ListItem[T Indentable[T]] struct {
	expanded    bool
	Value       T
	ParentModel *Model[T]
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
	Value() T
	ParentOptions() Options
	SetOptions(Options)
}

func sliceNFind[T Indentable[T]](cur []Indentable[T]) int {
	accum := 0
	for _, i := range cur {
		accum += 1
		accum += sliceNFind[T](i.Children())
	}
	return accum
}

func (i *ListItem[T]) Add(item Indentable[T], index int) {
	listItem := ListItem[T]{
		Value:       item.Value(),
		ParentModel: i.ParentModel,
	}
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
	if i.Value.ParentOptions().Expandable {
		i.expanded = v
	}
}

func NewItem[T Indentable[T]](item T, del list.ItemDelegate) ListItem[T] {
	li := ListItem[T]{
		Value:       item,
		ParentModel: &Model[T]{
		},
	}
	li.ParentModel.list = list.New([]list.Item{}, del, 200, 200)
	return li
}