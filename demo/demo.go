package demo

import (
	"fmt"
	"io"
	"strings"

	"git.tablet.sh/tablet/boba/datepicker"
	util "git.tablet.sh/tablet/boba/utilTypes"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#00e7e3")).MarginLeft(3).Blink(true)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(5)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(3).Foreground(lipgloss.Color("#affffe"))
	confirmTextStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#dddddd")).Margin(1, 0, 2, 4)
)

type DemoItem struct {
	text   string
	Result string
	model  *ModelContainer
}

type ModelContainer struct {
	value *tea.Model
}

func (i DemoItem) FilterValue() string { return i.text }

type itemDeleg struct{}

func (d itemDeleg) Height() int                               { return 1 }
func (d itemDeleg) Spacing() int                              { return 0 }
func (d itemDeleg) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDeleg) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(DemoItem)
	if !ok {
		return
	}
	str := fmt.Sprintf("[%d]. %s", index+1, i.text)
	fn := itemStyle.Render

	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("♣ " + strings.Join(s, " "))
		}
	}
	fmt.Fprint(w, fn(str))
}

type DemoModel struct {
	List        list.Model
	choice      string
	demoStarted bool
}

func (m DemoModel) Init() tea.Cmd {
	return nil
}

func (m DemoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	var anon func() = func() {
		i, _ := m.List.SelectedItem().(DemoItem)
		nope, cmd := (*(*i.model).value).Update(msg)
		*(m.List.SelectedItem().(DemoItem).model) = ModelContainer{
			value: &nope,
		}
		cmds = append(cmds, cmd)
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if !m.demoStarted {
				m.demoStarted = true
			} else {
				anon()
			}
		default:
			anon()
		}
	case util.GenResultMsg[string]:
		tea.Printf("\n---return value\n---\n [ %s ]", m.List.SelectedItem().(DemoItem).Result)
		return m, tea.Quit
	}
	if !m.demoStarted {
		m.List, cmd = m.List.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m DemoModel) View() string {
	if m.demoStarted {
		return (*(*m.List.SelectedItem().(DemoItem).model).value).View()
	} else if m.choice != "" {
		return confirmTextStyle.Render(fmt.Sprintf("demoing bubble : %s", m.choice))
	} else {
		return "\n" + m.List.View()
	}
}

func Setup() DemoModel {
	var modi tea.Model = datepicker.Initialize()
	items := []list.Item{
		DemoItem{
			text: "Date picker",
			model: &ModelContainer{
				value: &modi,
			},
			Result: "",
		},
	}
	l := list.New(items, itemDeleg{}, 20, 15)
	l.Title = "Select a Boba(TM) component to demo."
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle

	m := DemoModel{List: l, demoStarted: false}

	return m
}
