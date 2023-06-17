package recursivelist

import (
	// "fmt"

	"math"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ListItem[T Indentable[T]] struct {
	expanded    bool
	Value       *T
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
	Value() *T
	ParentOptions() Options
	SetOptions(Options)
	TotalBeneath() int
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
	v := item.Value()
	listItem := ListItem[T]{
		Value:       v,
		ParentModel: i.ParentModel,
	}
	// i.Children() = append(i.Children(), listItem)

	var top *T = item.Parent()
	accum := item.IndexWithinParent()
	for top != nil {
		iwi := (*top).IndexWithinParent()
		if iwi >= 0 {
			accum += iwi + 1
			slic := int(math.Min(
				float64(iwi),
				float64(len((*top).Children())),
			))
			p := (*top)
			for _, wee := range (p).Children()[0:slic] {
				accum += sliceNFind[T](wee.Children())
			}
		}

		top = (*top).Parent()
	}
	i.ParentModel.list.InsertItem(accum + index, listItem)
}

func (i ListItem[T]) Flatten() []T {
	return (*i.Value).Flatten()
}

func (i ListItem[T]) FilterValue() string {
	return (*i.Value).FilterValue()
}
func (i ListItem[T]) Expanded() bool {
	return i.expanded
}

func (i ListItem[T]) Point() *ListItem[T] {
	return &i
}

func (i *ListItem[T]) SetExpanded(v bool) tea.Cmd {
	if (*i.Value).ParentOptions().Expandable {
		i.expanded = v
	}
	if !v {
		for ii := 0; ii < (*i.Value).TotalBeneath(); ii++ {
			i.ParentModel.list.RemoveItem(ii)
		}
	}
	return nil
}

func NewItem[T Indentable[T]](item T, del list.ItemDelegate) ListItem[T] {
	li := ListItem[T]{
		Value:       &item,
		ParentModel: &Model[T]{
		},
		expanded: true,
	}
	li.ParentModel.list = list.New([]list.Item{}, del, 200, 200)
	return li
}