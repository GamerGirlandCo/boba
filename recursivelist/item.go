package recursivelist

import (
	// "fmt"

	"github.com/charmbracelet/bubbles/list"
)

type ListItem[T Indentable] struct {
	children    *[]ListItem[T]
	expanded    bool
	Value       T
	listoptions Options
	ParentModel *Model[T]
	Indentable
}

type Indentable interface {
	list.Item
	Lvl() int
	FilterValue() string
}

func (i ListItem[T]) newParentModel() list.Model {
	lis := make([]list.Item, 0)
	for _, inter := range *i.children {
		lis = append(lis, inter)
	}
	// fmt.Printf("PMPMPMPMPMPMPMPMPM %+v", *i.ParentModel.options)
	// fmt.Printf("lilililililililili %+v", lis)
	d := (*i.ParentModel).Delegate
	wi := (*i.ParentModel).options.Width
	hi := (*i.ParentModel).options.Height
	m := list.New(lis, d, wi, hi)
	return m
}

func (i ListItem[T]) Flatten() []list.Item {
	accum := make([]list.Item, 0)
	for _, ite := range *i.children {
		accum = append(accum, ite)
		accum = append(accum, ite.Flatten()...)
	}
	return accum
}

func (i ListItem[T]) FilterValue() string {
	return i.Value.FilterValue()
}
func (i ListItem[T]) Expanded() bool {
	return i.expanded
}

func (i ListItem[T]) ListOptions() Options {
	return i.listoptions
}

func (i *ListItem[T]) SetExpanded(v bool, pm Model[T]) {
	if pm.options.Expandable {
		i.expanded = v
	}
}

func (i *ListItem[T]) recurseAndExpand(pm Model[T]) {
	i.SetExpanded(true, pm)
	for _, i := range *i.children {
		i.recurseAndExpand(pm)
	}
}

func (i *ListItem[T]) AddItem(item T, index int) {
	listItem := ListItem[T]{
		Value:     item,
		ParentModel: i.ParentModel,
		children: &[]ListItem[T]{},
	}
	*i.children = append(*i.children, listItem)

	// fmt.Printf("lalalalala\n %+v\n", i.ParentModel.list)
	i.ParentModel.list.InsertItem(index, item)
}

func (i *ListItem[T]) AddMulti(items ...T) {
	leng := len(*i.children)
	for indy, ii := range items {
		i.AddItem(ii, leng+indy)
	}
}

func NewItem[T Indentable](item T, del list.ItemDelegate) ListItem[T] {
	li := ListItem[T]{
		Value:       item,
		children:    &[]ListItem[T]{},
		ParentModel: &Model[T]{},
	}
	li.ParentModel.list = list.New([]list.Item{}, del, 200, 200)
	return li
}
