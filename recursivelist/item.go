package recursivelist

import (
	// "fmt"
	"reflect"
	"strings"
	// "io"

	"git.tablet.sh/tablet/boba/types"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type ListItem[T list.Item] struct {
	Component   list.Model
	children    *[]ListItem[T]
	parent      *ListItem[T]
	expanded    bool
	Value       T
	listoptions Options
	ParentModel *Model[T]
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
	comp := i.newParentModel()
	listItem := ListItem[T]{
		Value:     item,
		Component: comp,
		parent: i,
		ParentModel: i.ParentModel,
		children: &[]ListItem[T]{},
	}
	*i.children = append(*i.children, listItem)

	// fmt.Printf("lalalalala %+v\n", i.Component)
	i.Component.InsertItem(index, item)
}

func (i *ListItem[T]) AddMulti(items ...T) {
	leng := len(*i.children)
	for indy, ii := range items {
		i.AddItem(ii, leng+indy)
	}
}

func (i ListItem[T]) findIndent(lv *int) int {
	whilevar := *i.parent
	ret := 0
	for whilevar.parent != nil {
		ret++
		whilevar = *whilevar.parent
	}
	return ret
}
func (i *ListItem[T]) setItemConfig(o ListOptions) {
	filds := reflect.TypeOf(i.ParentModel.options)
	values := reflect.ValueOf(i.ParentModel.options)

	for j := 0; j < filds.NumField(); j++ {
		f := filds.Field(j)
		v := values.Field(j)
		if veel, indo := types.FindField(o, f.Name); indo != -1 {
			v.Set(*veel)
		}
		curFack := reflect.TypeOf(i.Component)
		for k := 0; k < curFack.NumField(); k++ {
			if curFack.Field(k).IsExported() && curFack.Name() == f.Name {
				reflect.ValueOf(i.Component).Field(k).Set(v.Elem())
			}
		}
	}
	for _, item := range *i.children {
		item.setItemConfig(o)
	}
}

func (i *ListItem[T]) SetSize(w, h int) {
	i.Component.SetSize(w, h)
	for _, i := range *i.children {
		i.SetSize(w, h)
	}
}

func (i ListItem[T]) Init() tea.Cmd {
	return nil
}

func (i ListItem[T]) View() string {
	sb := strings.Builder{}
	var np int = 0
	for _, val := range *i.children {
		nesto := val.findIndent(&np) * 2
		indStyle := lg.NewStyle().Padding(0, 0, 0, nesto)
			sb.WriteString(indStyle.Render(
				val.Component.View(),
			))
			sb.WriteString("\n")
			sb.WriteString(val.View())
	}
	return sb.String()
}
func (i ListItem[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	i.Component, cmd = i.Component.Update(msg)
	cmds = append(cmds, cmd)
	return i, tea.Batch(cmds...)
}

func NewItem[T list.Item](item T, lid list.ItemDelegate) ListItem[T] {
	compy := list.New([]list.Item{}, lid, 200, 400)
	li := ListItem[T]{
		Value:       item,
		children:    &[]ListItem[T]{},
		ParentModel: &Model[T]{},
		Component: compy,
	}
	return li
}
