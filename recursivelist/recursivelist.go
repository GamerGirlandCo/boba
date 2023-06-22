package recursivelist

import (
	// "fmt"

	"log"
	"os"

	"git.tablet.sh/tablet/boba/utils"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
	// lg "github.com/charmbracelet/lipgloss"
)

type KeyMap struct {
	Expand key.Binding
	Choose key.Binding
}

var DefaultKeys KeyMap = KeyMap{
	Expand: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "expand/collapse item"),
	),
	Choose: key.NewBinding(
		key.WithKeys(" ", tea.KeyEnter.String()),
		key.WithHelp("<space>/↲", "choose item"),
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
	items   []ListItem[T]
	List    list.Model
	Options Options
}

// creates a new ListItem.
// note that this does NOT append it to
// m.items; you have to do that yourself
// with m.AddToRoot()
func (m *Model[T]) NewItem(item T) ListItem[T] {
	var childVar []ListItem[T]
	// var pari *ListItem[T] = nil

	for _, val := range item.GetChildren() {
		childVar = append(childVar, m.NewItem(val))
	}
	litem := &item
	expanded := true
	li := ListItem[T]{
		value:       litem,
		ParentModel: m,
		Children:    &childVar,
		expanded:    &expanded,
	}
	// 	if item.GetParent() != nil {
	// 		point := m.NewItem(*item.GetParent())
	// 		li.Parent = &point
	// 	}
	return li
}

func (m *Model[T]) SetSize(w, h int) {
	m.List.SetSize(w, h)
}

func (m *Model[T]) recurseAndExpand(i ListItem[T], currentState bool) tea.Cmd {
	var cmds []tea.Cmd
	cur := m.List.SelectedItem().(ListItem[T])
	iwith := m.List.Index()
	tots := cur.TotalBeneath() + 1
	if !currentState {
		for wtf := 0; wtf < tots; wtf++ {
			m.List.RemoveItem(wtf + iwith)
		}
	} else {
		toInsrt := cur.Flatten()
		for ma := range toInsrt {
			cmds = append(cmds, m.List.InsertItem(ma+iwith+1, toInsrt[ma]))
		}
	}
	i.SetExpanded(currentState)
	m.List.SetShowPagination(true)
	cmds = append(cmds, m.List.SetItem(m.List.Index(), i))
	cmd, _ := m.Flatten()
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}

func (m *Model[T]) AddToRoot(argument ...T) {
	var lips []ListItem[T]
	for _, mop := range argument {
		lips = append(lips, m.NewItem(mop))
	}
	m.items = lips
	_, t := m.Flatten()
	m.List.SetItems(t)
}

func (m *Model[T]) Flatten() (tea.Cmd, []list.Item) {
	accum := make([]list.Item, 0)
	for _, ite := range m.items {
		ite.ParentModel = m
		if *ite.expanded {
			for _, b := range ite.Flatten() {
				accum = append(accum, b)
			}
		} else {
			accum = append(accum, ite)
		}
	}
	lak := []tea.Cmd{
		m.List.SetItems(accum),
	}
	return tea.Batch(lak...), accum
}

func (m Model[T]) Init() tea.Cmd {
	cmd, toSet := m.Flatten()
	return tea.Batch(cmd, m.List.SetItems(toSet))
}

func (m Model[T]) View() string {
	lak := m.List.View()
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

			cmds = append(cmds, m.recurseAndExpand(blip, !*blip.expanded))
		case key.Matches(msg, m.Options.Keymap.Choose):
			return m, func() tea.Msg {
				reso := m.List.SelectedItem().(ListItem[T])
				result := utils.GenResultMsg{
					Res:       *reso.value,
					StringRep: (*reso.value).ToString(),
				}
				return result
			}
		}
	case tea.MouseMsg:
		log.Print("it is a mouse.", msg)
	}
	// cmds = append(cmds, m.List.SetItems(toset))
	nlm, cmd := m.List.Update(msg)
	m.List = nlm
	cmds = append(cmds, cmd)
	cmd, _ = m.Flatten()
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func New[T ItemWrapper[T]](items []T, delegate list.ItemDelegate, options Options) Model[T] {
	lis := make([]list.Item, 0)
	m := Model[T]{
		items:   []ListItem[T]{},
		Options: options,
	}
	if !m.Options.Expandable {
		m.Options.Keymap.Expand.SetEnabled(false)
	}
	for iii, it := range items {
		lis = append(lis, it)
		ni := m.NewItem(it)
		*ni.ParentModel = m
		m.items = append(m.items, ni)
		*m.items[iii].ParentModel = m
	}
	wpo, h, _ := term.GetSize(int(os.Stdout.Fd()))
	m.List = list.New(lis, delegate, wpo, h)
	m.List.Styles = list.DefaultStyles()
	m.List.SetFilteringEnabled(false)
	m.List.InfiniteScrolling = true
	m.List.AdditionalShortHelpKeys = func() []key.Binding {
		return utils.IterKeybindings(m.Options.Keymap)
	}
	m.Flatten()
	return m
}
