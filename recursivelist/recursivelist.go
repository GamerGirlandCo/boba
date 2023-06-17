package recursivelist

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
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

type Model[T list.Item] struct {
	options  *Options
	items    []ListItem[T]
	Delegate list.ItemDelegate
}

func (m *Model[T]) SetConfig(o Options) {
	m.options = &o
	for _, val := range m.items {
		val.setItemConfig(o.ListOptions)
	}
}

func (m *Model[T]) SetSize(w, h int) {
	for _, val := range m.items {
		val.SetSize(w, h)
	}
}

func (m *Model[T]) SetExpandable(v bool) {
	m.options.Expandable = v
	for _, i := range m.items {
		if !v {
			i.recurseAndExpand(*m)
		}
	}
}

func (m *Model[T]) Expandable() bool {
	return m.options.Expandable
}

func (m *Model[T]) SetItems(argument []ListItem[T]) {
	m.items = argument
}

func (i Model[T]) Init() tea.Cmd {
	// return tea.EnterAltScreen
	return nil
}

func (m Model[T]) View() string {
	sb := strings.Builder{}

	for _, it := range m.items {
		sb.WriteString(it.View())
		sb.WriteString("\n")
	}

	return sb.String()
}

func (m Model[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:

	}
	for _, ra := range m.items {
		nlm, cmd := ra.Component.Update(msg)
		ra.Component = nlm
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

type ItemInterface interface {
	Title() string
	FilterValue() string
}

func New[T list.Item](items []T, delegate list.ItemDelegate, width, height int) Model[T] {
	m := Model[T]{
		options: &Options{
			ClosedPrefix: ">",
			OpenPrefix:   "⌵",
			Width:        width,
			Height:       height,
		},
		Delegate: delegate,
		items:    []ListItem[T]{},
	}
	lis := make([]list.Item, 0)
	for iii, it := range items {
		lis = append(lis, it)
		ni := NewItem[T](it, delegate)
		*ni.ParentModel = m
		m.items = append(m.items, ni)
		lm := list.New(lis, delegate, width, height)
		(m.items[iii]).Component = lm
		lm.SetFilteringEnabled(false)
		m.items[iii].ParentModel = &m
	}
	return m
}
