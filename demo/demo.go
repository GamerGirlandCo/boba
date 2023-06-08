package demo

import (
	"fmt"
	"io"
	"strings"

	"git.tablet.sh/tablet/boba/datepicker"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	util "git.tablet.sh/tablet/boba/utilTypes"
)

var (
	titleStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#00e7e3")).MarginLeft(3).Blink(true)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(5)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(3).Foreground(lipgloss.Color("#affffe"))
	confirmTextStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#dddddd")).Margin(1, 0, 2, 4)
)

type DemoItem struct {
	text  string
	model tea.Model
	Result string
}

func (i DemoItem) FilterValue() string { return i.text }

type runMsg struct{}

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
	List     list.Model
	choice   string
	quitting bool
}

func (m DemoModel) Init() tea.Cmd {
	return nil
}

func (m DemoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			m.quitting = true
			i, ok := m.List.SelectedItem().(DemoItem)
			if ok {
				m.choice = i.text
				return i.model, nil
			} else {
				return m, tea.Quit
			}
		}
	case runMsg:
		return m, tea.ClearScreen
	case util.GenResultMsg:
		return m, tea.Quit
	}

	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m DemoModel) View() string {
	if m.choice != "" {
		return confirmTextStyle.Render(fmt.Sprintf("demoing bubble : %s", m.choice))
	} else {
		return "\n" + m.List.View()
	}
}

func Setup() DemoModel {
	items := []list.Item{
		DemoItem{
			text:  "Date picker",
			model: datepicker.Initialize(),
			Result: "",
		},
	}
	l := list.New(items, itemDeleg{}, 20, 15)
	l.Title = "Select a Boba(TM) component to demo."
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle

	m := DemoModel{List: l}

	return m
}
