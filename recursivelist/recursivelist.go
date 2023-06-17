package recursivelist

import (
	// "fmt"

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
	FilterintEnabled bool
	InfiniteScrolling bool
}

type Model[T Indentable[U], U list.Item] struct {
	options  *Options
	items    []ListItem[T,U]
	Delegate list.ItemDelegate
	list list.Model
}

func (m *Model[T, U]) SetSize(w, h int) {
	m.list.SetSize(w, h)
}

func (m *Model[T, U]) SetExpandable(v bool) {
	m.options.Expandable = v
	for _, i := range m.items {
		if !v {
			i.recurseAndExpand(*m)
		}
	}
}

func (m *Model[T, U]) Expandable() bool {
	return m.options.Expandable
}

func (m *Model[T, U]) SetItems(argument []ListItem[T, U]) {
	m.items = argument
}

func (m *Model[T, U]) Flatten() tea.Cmd {
	accum := make([]list.Item, 0)
	for _, ite := range m.items {
		accum = append(accum, ite)
		accum = append(accum, ite.Flatten()...)
	}
	return m.list.SetItems(accum)

}

func (i Model[T, U]) Init() tea.Cmd {
	// return tea.EnterAltScreen
	return nil
}

func (m Model[T, U]) View() string {
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

func (m Model[T, U]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	// case tea.KeyMsg:

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

func New[T Indentable[U], U list.Item](items []T, delegate list.ItemDelegate, width, height int) Model[T, U] {
	lis := make([]list.Item, 0)
	m := Model[T, U]{
		options: &Options{
			ClosedPrefix: ">",
			OpenPrefix:   "⌵",
			Width:        width,
			Height:       height,
		},
		Delegate: delegate,
		items:    []ListItem[T, U]{},
		list: list.New(lis, delegate, width, height),
	}
	m.list.Styles = list.DefaultStyles()
	m.list.SetFilteringEnabled(false)
	for iii, it := range items {
		lis = append(lis, it)
		ni := NewItem[T, U](it, delegate)
		*ni.ParentModel = m
		m.items = append(m.items, ni)
		// (m.items[iii]).Component = lm
		*m.items[iii].ParentModel = m
	}
	m.Flatten()
	return m
}
