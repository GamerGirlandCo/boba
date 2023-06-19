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

var DefaultOptions Options = Options{
	ClosedPrefix: ">",
	OpenPrefix:   "⌵",
	Width:        600,
	Height:       250,
	Expandable:   true,
	Keymap:       DefaultKeys,
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
	List     list.Model
	Options  Options
}

func (m *Model[T]) SetSize(w, h int) {
	m.List.SetSize(w, h)
}

func (m *Model[T]) recurseAndExpand(i ListItem[T]) tea.Cmd {
	var cmds []tea.Cmd
	for _, ee := range m.List.Items() {
		curLevel := (*ee.(ListItem[T]).value).Lvl()
		lastLvl := (*i.Value()).Lvl()
		if curLevel <= lastLvl {
			cur := ee.(ListItem[T])
			cur.Point().SetExpanded(!cur.expanded)
			iwith := cur.IndexWithinParent()
			tots := cur.TotalBeneath()
			if !cur.expanded {
				for ichrist := iwith; ichrist < iwith + tots + 1; ichrist++ {
					cur.ParentModel.List.RemoveItem(ichrist)
				}
			} else {
				toInsrt := cur.Flatten()
				for m := range toInsrt {
					cur.ParentModel.List.InsertItem(m + iwith, toInsrt[m])
				}
			}
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
		// m.List.SetItems([]List.Item{}),
		m.List.SetItems(accum),
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
	lak := m.List.View()
	// fmt.Println(lak)
	return lak
}

func (m Model[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.List.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Options.Keymap.Expand):
			blip := m.List.SelectedItem().(ListItem[T])
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
	nlm, cmd := m.List.Update(msg)
	m.List = nlm
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func New[T ItemWrapper[T]](items []T, delegate list.ItemDelegate, options Options) Model[T] {
	lis := make([]list.Item, 0)
	m := Model[T]{
		Delegate: delegate,
		items:    []ListItem[T]{},
		Options:  options,
	}
	for iii, it := range items {
		lis = append(lis, it)
		ni := NewItem(it, delegate, options)
		*ni.ParentModel = m
		m.items = append(m.items, ni)
		*m.items[iii].ParentModel = m
	}
	_, h, _ := term.GetSize(int(os.Stdout.Fd()))
	m.List = list.New(lis, delegate, 0, h)
	// m.List.Paginator = paginator.New()
	// m.List.Paginator.PerPage = 10
	m.List.Styles = list.DefaultStyles()
	m.List.SetFilteringEnabled(false)
	// m.List.InfiniteScrolling = true
	m.Flatten()
	return m
}
