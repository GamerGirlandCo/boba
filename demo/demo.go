package demo

import (
	"fmt"
	"io"
	"strings"

	"git.tablet.sh/tablet/boba/datepicker"
	"git.tablet.sh/tablet/boba/timepicker"
	util "git.tablet.sh/tablet/boba/utils"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	z "github.com/lrstanley/bubblezone"
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
			return selectedItemStyle.Render("🌸 " + strings.Join(s, " "))
		}
	} else {
	  str = "   " + str
	}
	fmt.Fprint(w, fn(str))
}

type DemoModel struct {
	List        list.Model
	choice      string
	choiceInd   int
	demoStarted bool
	items       []*tea.Model
}

func (m DemoModel) Init() tea.Cmd {
	return nil
}

func (m DemoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
	case util.GenResultMsg:
		tea.Printf("\n---\nreturn value\n---\n [ %s ]", msg.Res)
		var i DemoItem = m.List.SelectedItem().(DemoItem)
		*i.Result = msg.StringRep
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	}
	if !m.demoStarted {
		return updateChoices(m, msg)
	} else {
		return updateChosen(m, msg)
	}
}

func (m DemoModel) View() string {
	result := ""
	if r := m.List.SelectedItem().(DemoItem).Result; *r != "" {
		result = *r
	}
	if m.demoStarted {
		return z.Scan(
			lipgloss.JoinVertical(lipgloss.Center,
				confirmTextStyle.Render(fmt.Sprintf("demoing bubble: %s", m.choice)),
				chosenView(m),
				result,
			))
		// )
	} else {
		return choicesView(m)
	}
}

func updateChoices(m DemoModel, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.choiceInd = m.List.Cursor()
			m.demoStarted = true
			return m, nil
		default:
			taki, up := m.List.Update(msg)
			m.List = taki
			return m, up
		}
	}
	return m, nil
}


func updateChosen(m DemoModel, msg tea.Msg) (tea.Model, tea.Cmd) {
	tako, cmd := (*m.items[m.choiceInd]).Update(msg)
	*m.items[m.choiceInd] = tako
	return m, cmd
}

func choicesView(m DemoModel) string {
	return "\n" + m.List.View()
}

func chosenView(m DemoModel) string {
	return (*m.items[m.List.Cursor()]).View()
}

func Setup() DemoModel {
	titles := []string{
		"Date picker",
		"Time picker",
		"Time picker (with seconds)",
		"Recursive list",
	}
	items := make([]list.Item, len(titles))
	var modi []*tea.Model
	for i := range items {
		var dp tea.Model
		switch i {
		case 0:
			dp = datepicker.Initialize()
		case 1:
			dp = timepicker.Initialize(false)
		case 2:
			dp = timepicker.Initialize(true)
		case 3:
			dp = initRlistModel()
		}
		modi = append(modi, &dp)
		var minit string = ""
		items[i] = DemoItem{
			text:   titles[i],
			Result: &minit,
		}
	}

	l := list.New(items, itemDeleg{}, 20, 15)
	l.Title = "🧋 Select a Boba component to demo."
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle

	m := DemoModel{List: l, demoStarted: false, items: modi}

	return m
}
