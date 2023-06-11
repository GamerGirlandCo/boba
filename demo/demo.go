package demo

import (
	"fmt"
	"io"
	"strings"

	"git.tablet.sh/tablet/boba/datepicker"
	"git.tablet.sh/tablet/boba/timepicker"
	util "git.tablet.sh/tablet/boba/types"
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
	Result *string
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
		anon()

	case tea.MouseMsg:
		anon()
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if !m.demoStarted {
				m.demoStarted = true
				m.choice = m.List.SelectedItem().(DemoItem).text
			} else {
				anon()
			}
		default:
			anon()
		}
	case util.GenResultMsg[string]:
		tea.Printf("\n---return value\n---\n [ %s ]", msg.Res)
		var i DemoItem = m.List.SelectedItem().(DemoItem)
		*i.Result = msg.Res
		return m, tea.Quit
	}
	if !m.demoStarted {
		m.List, cmd = m.List.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m DemoModel) View() string {
	result := ""
	if r := m.List.SelectedItem().(DemoItem).Result; *r != "" {
		result = *r
	}
	if m.demoStarted {
		return lipgloss.JoinVertical(lipgloss.Center, 
			confirmTextStyle.Render(fmt.Sprintf("demoing bubble : %s", m.choice)),
			(*(*m.List.SelectedItem().(DemoItem).model).value).View(),
			result,
			)
	} else {
		return "\n" + m.List.View()
	}
}

func Setup() DemoModel {

	titles := []string{
		"Date picker",
		"Time picker",
		"Time picker (with seconds)",
	}
	items := make([]list.Item, len(titles))
	for i := range items {
		var modi tea.Model
		switch i {
		case 0:
			modi = datepicker.Initialize()
		case 1:
			modi = timepicker.Initialize(false)
		case 2:
			modi = timepicker.Initialize(true)
		}
		var minit string = ""
		items[i] = DemoItem{
			text: titles[i],
			model: &ModelContainer{
				value: &modi,
			},
			Result: &minit,
		}
	}

	l := list.New(items, itemDeleg{}, 20, 15)
	l.Title = "🌸 Select a Boba🧋 component to demo."
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle

	m := DemoModel{List: l, demoStarted: false}

	return m
}
