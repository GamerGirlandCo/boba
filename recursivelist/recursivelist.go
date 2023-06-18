package recursivelist

import (
	// "fmt"

	"log"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
	// lg "github.com/charmbracelet/lipgloss"
)

type KeyMap struct {
	list.KeyMap
	Expand key.Binding
	Choose key.Binding
	Quit   key.Binding
}

var DefaultKeys KeyMap = KeyMap{
	Expand: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "expand/collapse"),
	),
}

type Options struct {
	ClosedPrefix string
	OpenPrefix   string
	ListOptions  ListOptions
	Expandable   bool
	Width        int
	Height       int
	Keymap       KeyMap
}

func (o *Options) SetExpandable(v bool) {
	o.Expandable = v
}

type ListOptions struct {
	Keymap            list.KeyMap
	Styles            list.Styles
	Title             string
	FilteringEnabled  bool
	InfiniteScrolling bool
}

type Model[T ItemWrapper[T]] struct {
	items    []ListItem[T]
	Delegate list.ItemDelegate
	list     list.Model
	Options  Options
}

func (m *Model[T]) SetSize(w, h int) {
	m.list.SetSize(w, h)
}

func (m *Model[T]) recurseAndExpand(i ListItem[T]) tea.Cmd {
	var cmds []tea.Cmd
	for _, ee := range m.list.Items() {
		if (*ee.(ListItem[T]).value).Lvl() > i.Value().Lvl() {
			ee.(ListItem[T]).Point().SetExpanded(true)
		}
	}
	return tea.Batch(cmds...)
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
		// m.list.SetItems([]list.Item{}),
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
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Options.Keymap.Expand):
			blip := m.list.SelectedItem().(ListItem[T])
			blip.expanded = !blip.expanded
			cmds = append(cmds, m.recurseAndExpand(blip))
		}
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

func New[T ItemWrapper[T]](items []T, delegate list.ItemDelegate, width, height int, options Options) Model[T] {
	lis := make([]list.Item, 0)
	m := Model[T]{
		Delegate: delegate,
		items:    []ListItem[T]{},
		Options:  options,
	}
	for iii, it := range items {
		lis = append(lis, it)
		ni := NewItem(it, delegate)
		*ni.ParentModel = m
		m.items = append(m.items, ni)
		*m.items[iii].ParentModel = m
	}
	_, h, _ := term.GetSize(int(os.Stdout.Fd()))
	m.list = list.New(lis, delegate, width, h - 10)
	// m.list.Paginator = paginator.New()
	// m.list.Paginator.PerPage = 10
	m.list.Styles = list.DefaultStyles()
	m.list.SetFilteringEnabled(false)
	// m.list.InfiniteScrolling = true
	m.Flatten()
	return m
}
