package recursivelist

import (
	// "fmt"

	"github.com/charmbracelet/bubbles/list"
)

type ListItem[T Indentable[U], U list.Item] struct {
	children    *[]ListItem[T, U]
	expanded    bool
	Value       T
	listoptions Options
	ParentModel *Model[T, U]
	Indentable[U]
}

type Indentable[T list.Item] interface {
	list.Item
	Lvl() int
	FilterValue() string
	Children() []T
}


func (i ListItem[T, U]) Flatten() []list.Item {
	accum := make([]list.Item, 0)
	for _, ite := range *i.children {
		accum = append(accum, ite)
		accum = append(accum, ite.Flatten()...)
	}
	return accum
}

func (i ListItem[T, U]) FilterValue() string {
	return i.Value.FilterValue()
}
func (i ListItem[T, U]) Expanded() bool {
	return i.expanded
}

func (i ListItem[T, U]) ListOptions() Options {
	return i.listoptions
}

func (i *ListItem[T, U]) SetExpanded(v bool, pm Model[T, U]) {
	if pm.options.Expandable {
		i.expanded = v
	}
}

func (i *ListItem[T, U]) recurseAndExpand(pm Model[T, U]) {
	i.SetExpanded(true, pm)
	for _, i := range *i.children {
		i.recurseAndExpand(pm)
	}
}

func (i *ListItem[T, U]) AddItem(item T, index int) {
	listItem := ListItem[T, U]{
		Value:     item,
		ParentModel: i.ParentModel,
		children: &[]ListItem[T, U]{},
	}
	*i.children = append(*i.children, listItem)

	// fmt.Printf("lalalalala\n %+v\n", i.ParentModel.list)
	i.ParentModel.list.InsertItem(index, item)
}

func (i *ListItem[T, U]) AddMulti(items ...T) {
	leng := len(*i.children)
	for indy, ii := range items {
		i.AddItem(ii, leng+indy)
	}
}

func NewItem[T Indentable[U], U list.Item](item T, del list.ItemDelegate) ListItem[T, U] {
	li := ListItem[T, U]{
		Value:       item,
		children:    &[]ListItem[T, U]{},
		ParentModel: &Model[T, U]{},
	}
	li.ParentModel.list = list.New([]list.Item{}, del, 200, 200)
	return li
}
