package datepicker

import (
	"log"
	"os"
	"strconv"

	// "strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
			Foreground(lipgloss.Color("#baaec4")).
			Align(lipgloss.Center).
			Padding(0, 1, 1, 2)
)

type dateCell struct {
	date     time.Time
	selected bool
	blank    bool
}

func (d *dateCell) toggle(selected bool) {
	(*d).selected = selected
}

func (d dateCell) Render() string {
	mystr := strconv.FormatInt(int64(d.date.Day()), 10)
	if len(mystr) == 1 {
		mystr = " " + mystr
	}
	if d.selected {
		return selectedCellStyle.Render(mystr)
	}
	if d.blank {
		return blankCellStyle.Render("")
	}
	return cellStyle.Render(mystr)
}

type Model struct {
	loaded       bool
	internalGrid [][]dateCell
	cursorY      int
	cursorX      int
	anchor       time.Time
}

func (m Model) FindIndex(fn Predicate[dateCell]) [][]int {
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
func (m *Model) wrapNav(key string) {
	var blah []int = []int{m.cursorY, m.cursorX}
	prevo := m.internalGrid[blah[0]][blah[1]]
	switch key {
	case "left":
		(*m).cursorX--
		if (*m).cursorX < 0 {
			log.Print("y < 0")
			(*m).cursorX = 6
			if (*m).cursorY > 0 {
				(*m).cursorY--
			}
		}
		diffo := m.internalGrid[m.cursorY][m.cursorX].date
		if diffo.Month() < prevo.date.Month() {
			(*m).updateAnchor((cal.MonthStart((*m).anchor.AddDate(0, 0, -1))), diffo)
		}
	case "right":
		(*m).cursorX++
		if (*m).cursorX > 6 {
			(*m).cursorX = 0
			if (*m).cursorY < 6-1 {
				(*m).cursorY++
			}
		}
		diffo := m.internalGrid[m.cursorY][m.cursorX].date
		if diffo.Month() > prevo.date.Month() {
			(*m).updateAnchor((cal.MonthStart((*m).anchor.AddDate(0, 0, 1))), diffo)
		}
	case "up":
		(*m).cursorY--
		if (*m).cursorY < 0 {
			(*m).cursorY = 6 - 1
		}
		diffo := m.internalGrid[m.cursorY][m.cursorX].date
		if diffo.Month() < prevo.date.Month() {
			(*m).updateAnchor((cal.MonthStart((*m).anchor.AddDate(0, 0, -7))), diffo)
		}
	case "down":
		(*m).cursorY++
		if (*m).cursorY > 6-1 {
			(*m).cursorY = 6 - 1
		}
		diffo := m.internalGrid[m.cursorY][m.cursorX].date
		if diffo.Month() > prevo.date.Month() {
			(*m).updateAnchor(cal.MonthStart((*m).anchor.AddDate(0, 0, 7)), diffo)
		}

	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.internalGrid[m.cursorY][m.cursorX].Select(false)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "down":
			m.wrapNav(msg.String())
		case tea.KeyLeft.String(), tea.KeyRight.String(), "left", "right":
			m.wrapNav(msg.String())
		case tea.KeyPgDown.String():
			addo := m.anchor.AddDate(0, 1, 0)
			m.updateAnchor(cal.MonthStart(addo), addo)
		case tea.KeyPgUp.String():
			addo := m.anchor.AddDate(0, -1, 0)
			m.updateAnchor(cal.MonthStart(addo), addo)
		}
		inGrid, _, _ := makeMatrix(m.internalGrid[m.cursorY][m.cursorX].date, m.cursorY, m.cursorX)
		m.internalGrid = inGrid
		found := m.FindIndex(func(dd dateCell) bool {
			return dd.selected
		})
		log.Printf("the length of findex == %d", len(found))
		for l, val := range found {
			log.Printf("found @ %d -- [%d][%d]", l, val[0], val[1])
		}
	}
	inGrid, _, _ := makeMatrix(m.internalGrid[m.cursorY][m.cursorX].date, m.cursorY, m.cursorX)
	m.internalGrid = inGrid
	// m.selectedDate = m.internalGrid[m.cursorY][m.cursorX].date
	return m, nil
}

func (m Model) View() string {
	s := ""
	axisY := make([]string, 0)
	wo, _, _ := term.GetSize(int(os.Stdout.Fd()))
	axisY = append(axisY, monthStyle.Render(lipgloss.PlaceHorizontal(wo, lipgloss.Center, m.internalGrid[0][6].date.Format("January 2006"))))

	header := make([]string, 0)
	for i := 0; i < 7; i++ {
		header = append(header, weekStyle.Render(m.internalGrid[0][i].date.Format("Mon")))
	}
	axisY = append(axisY, lipgloss.JoinHorizontal(lipgloss.Center, header...))
	for i := 0; i < len(m.internalGrid); i++ {
		axisY = append(axisY, m.renderWeek(i))
	}

	s += lipgloss.JoinVertical(lipgloss.Center, axisY...)
	return s
}

func (m Model) renderWeek(index int) string {
	longlong := make([]string, 1)
	for i := 0; i < len(m.internalGrid[index]); i++ {
		longlong = append(longlong, m.internalGrid[index][i].Render())
	}
	return lipgloss.JoinHorizontal(lipgloss.Center, longlong...)
}

func (m *Model) updateAnchor(argo time.Time, now time.Time) {
	m.anchor = argo
	yy, xx := getDefaultMatrix(now)
	ingrid, _, _ := makeMatrix(now, yy, xx)
	m.internalGrid = ingrid
	m.cursorY = yy
	m.cursorX = xx
}

func getDefaultMatrix(cur time.Time) (int, int) {
	som := cal.MonthStart(cur)
	_, fd := isoweek.FromDate(som.Year(), som.Month(), som.Day())
	fw := isoweek.StartTime(som.Year(), fd, time.Local).AddDate(0, 0, -1)

	// eom := cal.MonthEnd(cur)
	// _, ld := isoweek.FromDate(eom.Year(), eom.Month(), eom.Day())
	// lw := isoweek.StartTime(eom.Year(), ld, time.Local).AddDate(0, 0, 6)

	diffB := cur.Sub(fw)
	// diffA := lw.Sub(cur)

	modulo := (int(diffB.Hours()) / 24) % 7
	div := int((int(diffB.Hours()) / 24) / 7)

	return div, modulo
}

func makeMatrix(sel time.Time, ya int, xa int) ([][]dateCell, int, int) {
	g := make([][]dateCell, 6)
	startOfMonth := cal.MonthStart(sel)
	_, week := startOfMonth.ISOWeek()
	firstWeek := isoweek.StartTime(startOfMonth.Year(), week, time.Local).AddDate(0, 0, -1)
	// _, curWeek := time.Now().ISOWeek()
	// thisWeek := isoweek.StartTime(date.Year(), curWeek, time.Local)
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
				log.Printf("x and y: [%d][%d]", y, x)
				intY = x
				intX = y
				selBool = true
			} else if /* {
				log.Print("ya and xa == y and x")
				log.Printf("x and y: [%d][%d]", y, x)
				selBool = true
			}  else if */myDay.Month() < sel.Month() || myDay.Month() > sel.Month() {
				blankBool = true
				selBool = false
			}
			g[y][x] = dateCell{
				selected: selBool,
				date:     myDay,
				blank:    blankBool,
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
	rlnw := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)

	my, mx := getDefaultMatrix(rlnw)
	inGrid, y, x := makeMatrix(rlnw, my, mx)
	inGrid[my][mx].Select(true)
	inGrid[y][x].Select(true)
	startOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local)
	meep := Model{
		internalGrid: inGrid,
		cursorY:      y,
		cursorX:      x,
		anchor:       startOfMonth,
		loaded:       true,
	}
	return meep
}
