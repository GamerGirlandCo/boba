package datepicker

import (
	"fmt"
	"log"
	"os"
	"strconv"

	// "strings"
	"time"

	util "git.tablet.sh/tablet/boba/types"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	z "github.com/lrstanley/bubblezone"
	cal "github.com/rickar/cal/v2"
	"github.com/snabb/isoweek"
	"golang.org/x/term"
)

var (
	roundBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
	}

	cellStyle = lipgloss.NewStyle().
			Border(roundBorder, true).
			BorderForeground(lipgloss.Color("#5c5e9f")).
			Padding(0, 1).
			Foreground(lipgloss.Color("#ffffff")).
			Align(lipgloss.Center)

	monthStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#7542dc")).
			Foreground(lipgloss.Color("#ffffff")).
			Bold(true).
			Padding(2).
			Margin(2)
	selectedCellStyle = lipgloss.NewStyle().
				Border(roundBorder, true).
				BorderForeground(lipgloss.Color("#ff4884")).
				Padding(0, 1).
				Foreground(lipgloss.Color("#000000")).
				Background(lipgloss.Color("#ff4884")).
				Align(lipgloss.Center)
	blankCellStyle = lipgloss.NewStyle().
			Padding(1, 3).
			Align(lipgloss.Center)
	weekStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderTop(false).
			BorderBottom(true).
			BorderLeft(false).
			BorderRight(false).
			Foreground(lipgloss.Color("#baaec4")).
			Align(lipgloss.Center).
			Padding(0, 1, 0, 2)
)

type keyMap struct {
	Up          key.Binding
	Down        key.Binding
	Left        key.Binding
	Right       key.Binding
	PrevMonth   key.Binding
	NextMonth   key.Binding
	StartOfWeek key.Binding
	EndOfWeek   key.Binding
	Help        key.Binding
	Choose      key.Binding
	Quit        key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Choose, k.Help}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right, k.Choose, k.PrevMonth, k.NextMonth, k.StartOfWeek, k.EndOfWeek, k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	PrevMonth: key.NewBinding(
		key.WithKeys(tea.KeyPgUp.String(), "shift+up"),
		key.WithHelp("PgUp/Shift+↑", "go to prev month"),
	),
	NextMonth: key.NewBinding(
		key.WithKeys(tea.KeyPgDown.String(), "shift+down"),
		key.WithHelp("PgDn/Shift+↓", "go to next month"),
	),
	StartOfWeek: key.NewBinding(
		key.WithKeys("shift+left", "0"),
		key.WithHelp("Shift+←/0", "move to beginning of current week"),
	),
	EndOfWeek: key.NewBinding(
		key.WithKeys("shift+right", "$"),
		key.WithHelp("Shift+→/$", "move to end of current week"),
	),
	Help: key.NewBinding(
		key.WithKeys(tea.KeyCtrlQuestionMark.String()),
		key.WithHelp("?", "show help"),
	),
	Choose: key.NewBinding(
		key.WithKeys(tea.KeySpace.String(), tea.KeyEnter.String()),
		key.WithHelp("Enter/Space", "finalize selection"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("Ctrl+c/q", "cancel"),
	),
}

type dateCell struct {
	date     time.Time
	selected bool
	blank    bool
	id string
}

func (d *dateCell) toggle(selected bool) {
	(*d).selected = selected
}

func (d dateCell) Render() string {
	mystr := strconv.FormatInt(int64(d.date.Day()), 10)
	// id := fmt.Sprintf("-%d-%s-%d", d.date.Year(), d.date.Format("01"), d.date.Day())
	if len(mystr) == 1 {
		mystr = " " + mystr
	}
	if d.selected {
		return z.Mark(d.id, selectedCellStyle.Render(mystr))
	}
	if d.blank {
		return z.Mark(d.id, blankCellStyle.Render(""))
	}
	return z.Mark(d.id, cellStyle.Render(mystr))
}

type Model struct {
	loaded       bool
	internalGrid [][]dateCell
	cursorY      int
	cursorX      int
	anchor       time.Time
	keys         keyMap
	help         help.Model
	value        time.Time
}

func (m Model) FindIndex(fn util.Predicate[dateCell]) [][]int {
	ret := make([][]int, 0)
	for x := 0; x < len(m.internalGrid)-1; x++ {
		for y := 0; y < len(m.internalGrid[x])-1; y++ {
			if fn(m.internalGrid[x][y]) {
				ret = append(ret, []int{x, y})
			}
		}
	}
	return ret
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.internalGrid[m.cursorY][m.cursorX].toggle(false)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width
	case tea.KeyMsg:
		var blah []int = []int{m.cursorY, m.cursorX}
		prevo := m.internalGrid[blah[0]][blah[1]]
		switch {
		// case key.Matches(msg, m.keys.Quit):
		// 	return m, tea.Quit
		case key.Matches(msg, m.keys.Down):
			m.cursorY++
			if m.cursorY > 6-1 {
				m.cursorY = 6 - 1
			}
			diffo := m.internalGrid[m.cursorY][m.cursorX].date
			if diffo.Month() > prevo.date.Month() {
				m.updateAnchor(cal.MonthStart(m.anchor.AddDate(0, 0, 7)), diffo)
			}
		case key.Matches(msg, m.keys.Up):
			m.cursorY--
			if m.cursorY < 0 {
				m.cursorY = 6 - 1
			}
			diffo := m.internalGrid[m.cursorY][m.cursorX].date
			if diffo.Month() < prevo.date.Month() {
				m.updateAnchor((cal.MonthStart(m.anchor.AddDate(0, 0, -7))), diffo)
			}
		case key.Matches(msg, m.keys.Left):
			m.cursorX--
			if m.cursorX < 0 {
				log.Print("y < 0")
				m.cursorX = 6
				if m.cursorY > 0 {
					m.cursorY--
				}
			}
			diffo := m.internalGrid[m.cursorY][m.cursorX].date
			if diffo.Month() < prevo.date.Month() {
				m.updateAnchor((cal.MonthStart(m.anchor.AddDate(0, 0, -1))), diffo)
			}
		case key.Matches(msg, m.keys.Right):
			m.cursorX++
			if m.cursorX > 6 {
				m.cursorX = 0
				if m.cursorY < 6-1 {
					m.cursorY++
				}
			}
			diffo := m.internalGrid[m.cursorY][m.cursorX].date
			if diffo.Month() > prevo.date.Month() {
				m.updateAnchor((cal.MonthStart(m.anchor.AddDate(0, 0, 1))), diffo)
			}
		case key.Matches(msg, m.keys.NextMonth):
			addo := m.anchor.AddDate(0, 1, 0)
			m.updateAnchor(cal.MonthStart(addo), addo)
		case key.Matches(msg, m.keys.PrevMonth):
			addo := m.anchor.AddDate(0, -1, 0)
			m.updateAnchor(cal.MonthStart(addo), addo)
		case key.Matches(msg, m.keys.StartOfWeek):
			m.cursorX = 0
			diffo := m.internalGrid[m.cursorY][m.cursorX].date
			if diffo.Month() < prevo.date.Month() {
				m.updateAnchor((cal.MonthStart(m.anchor.AddDate(0, 0, 1))), diffo)
			}
		case key.Matches(msg, m.keys.EndOfWeek):
			m.cursorX = 6
			diffo := m.internalGrid[m.cursorY][m.cursorX].date
			if diffo.Month() > prevo.date.Month() {
				m.updateAnchor((cal.MonthStart(m.anchor.AddDate(0, 0, 1))), diffo)
			}
		case key.Matches(msg, m.keys.Choose):
			m.value = m.internalGrid[m.cursorY][m.cursorX].date
			return m, func() tea.Msg {
				Res := m.value.String()
				return util.GenResultMsg[string]{
					Res: Res,
				}
			}
		}
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseLeft:
			for ia, a := range m.internalGrid {
				for ib, b := range a {
					if z.Get(b.id).InBounds(msg) {
						m.cursorY = ia
						m.cursorX = ib
					}
				}
			}
		}
	}
	inGrid, _, _ := makeMatrix(m.internalGrid[m.cursorY][m.cursorX].date, m.cursorY, m.cursorX)
	m.internalGrid = inGrid
	return m, nil
}

func (m Model) View() string {
	cal, _, _ := m.calendar()
	return lipgloss.JoinHorizontal(0.33, cal, m.help.View(m.keys))
}


// function to render out a calendar from the Mode 
// @receiver m 
// @return string 
// @return int 
// @return int 
func (m Model) calendar() (string, int, int) {
	s := ""
	otherRet := 0
	otherOtherRet := ""
	axisY := make([]string, 0)
	wo, _, _ := term.GetSize(int(os.Stdout.Fd()))
	axisY = append(axisY, z.Mark("title", monthStyle.Render(lipgloss.PlaceHorizontal(int(wo/2), lipgloss.Center, m.internalGrid[0][6].date.Format("January 2006")))))

	header := make([]string, 0)
	for i := 0; i < 7; i++ {
		kak := m.internalGrid[0][i].date
		header = append(header, z.Mark(fmt.Sprintf("wd-%s", kak.Format("Mon")), weekStyle.Render(kak.Format("Mon"))))
	}
	axisY = append(axisY, lipgloss.JoinHorizontal(lipgloss.Center, header...))
	otherRet = lipgloss.Height(lipgloss.JoinVertical(lipgloss.Center, axisY...))
	for i := 0; i < len(m.internalGrid); i++ {
		axisY = append(axisY, m.renderWeek(i))
		if i == 0 {
			otherOtherRet = m.renderWeek(i)
		}
	}

	s += lipgloss.JoinVertical(lipgloss.Center, axisY...)
	return z.Mark("fullcalendar", s), otherRet, lipgloss.Width(otherOtherRet)
}

func (m Model) renderWeek(index int) string {
	longlong := make([]string, 1)
	for i := 0; i < len(m.internalGrid[index]); i++ {
		longlong = append(longlong, m.internalGrid[index][i].Render())
	}
	y, w := m.internalGrid[index][0].date.ISOWeek()
	return z.Mark(fmt.Sprintf("w-%d-%d", y, w), lipgloss.JoinHorizontal(lipgloss.Center, longlong...))
}

func (m *Model) updateAnchor(argo time.Time, now time.Time) {
	m.anchor = argo
	yy, xx := getDefaultMatrix(now)
	ingrid, _, _ := makeMatrix(now, yy, xx)
	m.internalGrid = ingrid
	m.cursorY = yy
	m.cursorX = xx
}

func (m Model) Value() time.Time {
	return m.value
}

func getDefaultMatrix(cur time.Time) (int, int) {
	som := cal.MonthStart(cur)
	_, fd := isoweek.FromDate(som.Year(), som.Month(), som.Day())
	fw := isoweek.StartTime(som.Year(), fd, time.Local).AddDate(0, 0, -1)
	diffB := cur.Sub(fw)

	modulo := (int(diffB.Hours()) / 24) % 7
	div := int((int(diffB.Hours()) / 24) / 7)

	return div, modulo
}

func makeMatrix(sel time.Time, ya int, xa int) ([][]dateCell, int, int) {
	g := make([][]dateCell, 6)
	startOfMonth := cal.MonthStart(sel)
	_, week := startOfMonth.ISOWeek()
	firstWeek := isoweek.StartTime(startOfMonth.Year(), week, time.Local).AddDate(0, 0, -1)
	myDay := firstWeek

	var intY int
	var intX int
	for y := 0; y < 6; y++ {
		g[y] = make([]dateCell, 7)
		for x := 0; x < 7; x++ {
			var selBool bool = false
			var blankBool = false
			if (ya == x && xa == ya) && myDay.Equal(sel) {
				selBool = true
				intY = ya
				intX = xa
			} else if ya == y && xa == x {
				intY = x
				intX = y
				selBool = true
			} else if myDay.Month() < sel.Month() || myDay.Month() > sel.Month() {
				blankBool = true
				selBool = false
			}
			g[y][x] = dateCell{
				selected: selBool,
				date:     myDay,
				blank:    blankBool,
				id: z.NewPrefix(),
			}
			myDay = myDay.AddDate(0, 0, 1)
		}
	}

	// printGrid(g)
	return g, intX, intY
}

func logPos(m Model) {
	log.Printf("current!!! -- [%d][%d]", m.cursorY, m.cursorX)
}

func printGrid(d [][]dateCell) {
	for x, som := range d {
		s := ""
		for y := range som {
			s += d[x][y].date.Format("02 ")
		}
		log.Print(s)
	}
	log.Println(d[1][6].date.Format("January 2006"))
}

func Initialize() Model {
	z.NewGlobal()
	rlnw := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)

	my, mx := getDefaultMatrix(rlnw)
	inGrid, y, x := makeMatrix(rlnw, my, mx)
	inGrid[my][mx].toggle(true)
	inGrid[y][x].toggle(true)
	startOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local)
	hell := help.New()
	hell.ShowAll = true
	meep := Model{
		internalGrid: inGrid,
		cursorY:      y,
		cursorX:      x,
		anchor:       startOfMonth,
		loaded:       true,
		help:         hell,
		keys:         keys,
	}
	return meep
}
