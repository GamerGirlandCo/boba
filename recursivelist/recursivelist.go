package recursivelist

import (
	// "fmt"

	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	// lg "github.com/charmbracelet/lipgloss"
)

type Options struct {
	ClosedPrefix string
	OpenPrefix   string
	ListOptions  ListOptions
	Expandable   bool
	Width        int
	Height       int
}

type ListOptions struct {
	Keymap            list.KeyMap
	Styles            list.Styles
	Title             string
	FilterintEnabled  bool
	InfiniteScrolling bool
}

type Model[T Indentable[U], U list.Item] struct {
	Options  *Options
	items    []ListItem[T]
	Delegate list.ItemDelegate
	list     list.Model
}

func (m *Model[T, U]) SetSize(w, h int) {
	m.list.SetSize(w, h)
}

func (m *Model[T]) SetExpandable(v bool) {
	m.Options.Expandable = v
	for _, i := range m.items {
		if !v {
			m.recurseAndExpand(*m, i)
		}
	}
}

func (m *Model[T]) recurseAndExpand(pm Model[T], i ListItem[T]) {
	for _, ee := range m.list.Items() {
		if ee.(T).Lvl() > i.Value.Lvl() {
			ee.(ListItem[T]).Point().SetExpanded(true, *m)
		}
	}
}

func (m *Model[T]) NewItem(item T, del list.ItemDelegate) ListItem[T] {
	li := ListItem[T]{
		Value:       item,
		ParentModel: m,
	}
	li.ParentModel.list = list.New([]list.Item{}, del, 200, 200)
	return li
}

func (m *Model[T]) Expandable() bool {
	return m.Options.Expandable
}

func (m *Model[T]) SetItems(argument []ListItem[T]) {
	for _, mop := range argument {
		mop.ParentModel = m
		m.items = append(m.items, mop)
		in := len(m.items) - 1
		(*m).items[in] = mop
	}
	m.Flatten()
}

func (m *Model[T]) Flatten() tea.Cmd {
	accum := make([]list.Item, 0)
	for _, ite := range m.items {
		ite.ParentModel = m
		for _, b := range ite.Flatten() {
			accum = append(accum, b)
		}
	}
	lak := []tea.Cmd{
		m.list.SetItems([]list.Item{}),
		m.list.SetItems(accum),
	}
	return tea.Batch(lak...)
}

func (i Model[T]) Init() tea.Cmd {
	// return tea.EnterAltScreen
	return tea.Sequence(tea.EnterAltScreen, i.Flatten())
}

func (m Model[T]) View() string {
	// sb := strings.Builder{}
	// var np int = 0
	// for _, val := range *i.children {
	// 	nesto := val.findIndent(&np) * 2
	// 	indStyle := lg.NewStyle().Padding(0, 0, 0, nesto)

	// 		sb.WriteString("\n")
	// 		sb.WriteString(val.View())
	// }
	lak := m.list.View()
	// fmt.Println(lak)
	return lak
}

func (m Model[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	// case tea.KeyMsg:
	case tea.MouseMsg:
		log.Print("it is a mouse.", msg)
	}
	// for _, ra := range m.items {
	// nlm, cmd := ra.Component.Update(msg)
	// ra.Component = nlm
	// }
	cmds = append(cmds, m.Flatten())
	nlm, cmd := m.list.Update(msg)
	m.list = nlm
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}



func New[T Indentable[T]](items []T, delegate list.ItemDelegate, width, height int) Model[T] {
	lis := make([]list.Item, 0)
	m := Model[T]{
		Options: &Options{
			ClosedPrefix: ">",
			OpenPrefix:   "⌵",
			Width:        width,
			Height:       height,
		},
		Delegate: delegate,
		items:    []ListItem[T]{},
	}
	m.list.Styles = list.DefaultStyles()
	m.list.SetFilteringEnabled(false)
	for iii, it := range items {
		lis = append(lis, it)
		ni := m.NewItem(it, delegate)
		*ni.ParentModel = m
		m.items = append(m.items, ni)
		*m.items[iii].ParentModel = m
	}
	m.list = list.New(lis, delegate, 0, 0)
	m.Flatten()
	return m
}
