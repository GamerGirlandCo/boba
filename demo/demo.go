package demo

import (
	"fmt"
	"io"
	"strings"

	"git.tablet.sh/tablet/boba/datepicker"
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

type item struct {
	text string
	model tea.Model
}


func (i item) FilterValue() string { return i.text }

type runMsg struct{}

type itemDeleg struct{}

func (d itemDeleg) Height() int { return 1 }
func (d itemDeleg) Spacing() int { return 0 }
func (d itemDeleg) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDeleg) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
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

type model struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if !m.quitting {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.list.SetWidth(msg.Width)
			return m, nil
	
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "ctrl+c":
				return m, tea.Quit
			case "enter":
				m.quitting = true
				i, ok := m.list.SelectedItem().(item)
				if ok {
					m.choice = i.text
					return i.model, nil
				} else {
					return m, tea.Quit 
				}
			}
		case runMsg:
			return m, tea.ClearScreen
		}
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	if m.choice != "" {
		return confirmTextStyle.Render(fmt.Sprintf("demoing bubble : %s", m.choice))
	}
	return "\n" + m.list.View()
}

func Setup() model {
	items := []list.Item{
		item{
			text: "Date picker",
			model: datepicker.Initialize(),
		},
	}
	l := list.New(items, itemDeleg{}, 20, 15)
	l.Title = "Select a Boba(TM) component to demo."
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	
	m := model{list: l}

	return m
}