package timepicker

import (
	"fmt"
	"log"

	// "math"

	// "math"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

const (
	// upArrowSmall   string = ""
	// downArrowSmall string = "▔▀█▀▔"
	// upArrowWide    string = "▁▄▅▆▅▄▁"
	// downArrowWide  string = "▔▀▅▃▅▀▔"
	upArrowWide    string = "▁▄▆▄▁"
	downArrowWide  string = "▔▀█▀▔"
	upArrowSmall   string = "⌃⌃⌃"
	downArrowSmall string = "⌄⌄⌄"
	perMinute      int    = 60
	perHour        int    = 60
	perHalfDay     int    = 12
	perDay         int    = 24
)

type Styles struct {
	OuterBorder         lipgloss.Style
	SelectedOuterBorder lipgloss.Style
	Value               lipgloss.Style
	SelectedValue       lipgloss.Style
}

var defBorder = lipgloss.Border{
	Top:         "-",
	Bottom:      "-",
	Left:        "┆",
	Right:       "┆",
	TopLeft:     "+",
	TopRight:    "+",
	BottomLeft:  "+",
	BottomRight: "+",
}

type keyMap struct {
	NextField key.Binding
	PrevField key.Binding
	TickUp    key.Binding
	TickDown  key.Binding
	Choose    key.Binding
	First     key.Binding
	Last      key.Binding
}

var keys = keyMap{
	TickUp: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "value - 1")),
	TickDown: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "value + 1"),
	),
	PrevField: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "prev field"),
	),
	NextField: key.NewBinding(
		key.WithKeys(tea.KeyTab.String()),
		key.WithHelp("tab", "next field"),
	),
}

type ticky struct {
	val   int
	color string
}

var headers = []string{
	"Hour",
	"Minute",
	"Second",
}

type Model struct {
	// TwelveHr bool
	Seconds  bool
	Styles   Styles
	value    time.Time
	width    int
	selected int
	keys     keyMap
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	thirdField := 1
	if m.Seconds {
		thirdField = 3
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case tea.KeyMsg:
		switch {
		// case "q", "ctrl+c":
		// 	return m, tea.Quit
		case key.Matches(msg, m.keys.NextField):
			m.selected = minMax(0, thirdField, m.selected, 1, true, false)
		case key.Matches(msg, m.keys.PrevField):
			m.selected = minMax(0, thirdField, m.selected, -1, true, false)
		case key.Matches(msg, m.keys.TickUp):
			m.incOrDec(-1)
		case key.Matches(msg, m.keys.TickDown):
			m.incOrDec(1)
		}
	}

	return m, nil
}

func (m Model) View() string {
	no := 2
	if m.Seconds {
		no = 3
	}
	final := make([]string, no)
	for a := range final {
		tickers := genTicks(a, m)
		str := make([]string, 7)
		for b := range tickers {
			var sty lipgloss.Style
			if b == 3 {
				sty = m.Styles.SelectedValue
			} else {
				sty = m.Styles.Value.Copy().
					Foreground(lipgloss.Color(tickers[b].color))
				/* .Width(int(m.width / 3)) */
			}

			str[b] = sty.Render(fmt.Sprintf("%02d", tickers[b].val))
			if b == 6 {

				jd := lipgloss.JoinVertical(lipgloss.Top, str...)
				bs := m.Styles.OuterBorder
				if a == m.selected {
					bs = m.Styles.SelectedOuterBorder
				}
				final[a] = lipgloss.JoinVertical(
					lipgloss.Center, headers[a],
					upArrowWide, bs.Render(jd), downArrowWide)
			}
		}
	}

	return lipgloss.PlaceHorizontal(m.width, lipgloss.Center, lipgloss.JoinHorizontal(lipgloss.Center, final...))
}

func (m *Model) incOrDec(dir int) {
	base := m.value
	var sec, min, hr int = base.Second(), base.Minute(), base.Hour()
	switch m.selected {
	case 0:
		hr += dir
	case 1:
		min += dir
	case 2:
		sec += dir
	}
	delta := time.Date(base.Year(), base.Month(), base.Day(),
		hr, min, sec, 0, time.Local,
	)
	diff := delta.Sub(base)
	m.value = m.value.Add(diff)
}

func genTicks(t int, m Model) []ticky {
	ret := make([]ticky, 7)
	var wrapArg int
	var which int
	switch t {
	case 0:
		/* if m.TwelveHr {
			wrapArg = perHalfDay
			which = m.value.Hour() % 12
		} else */{
			wrapArg = perDay
			which = m.value.Hour()
		}
	case 1:
		wrapArg = perHour
		which = m.value.Minute()
	case 2:
		wrapArg = perMinute
		which = m.value.Second()
	}
	// j := 6
	grad := genGradient()
	for i := 0; i < 7; i++ {
		var sign int = -1
		// sub := 0
		var inten int
		if i < 3 {
			inten = i
		} else {
			// sub = 1
			inten = 7 - i - 1
		}
		// fmt.Println(inten)
		ret[i] = ticky{
			val: minMax(0, wrapArg,
				which+i,
				sign, true, true),
			color: grad[inten],
			// intensity: inten,
		}
	}
	return ret
}

func genGradient() []string {
	ograd := make([]colorful.Color, 4)
	x0, _ := colorful.Hex("#333333")
	x1, _ := colorful.Hex("#eeeeee")
	ret := make([]string, len(ograd))
	for i := range ograd {
		tooMuch := (float64(i) / 2)
		ograd[i] = x0.BlendHcl(x1, tooMuch)
		ret[i] = ograd[i].Hex()
	}
	log.Println(ret)
	return ret
}

func minMax(min, max, cur, dir int, wrapAround bool, inclusive bool) int {
	ret := cur
	secondsub := 0
	if dir < 0 {
		ret--
	} else if dir > 0 {
		ret++
	}
	if inclusive {
		secondsub = 1
	}
	// ret += dir
	if !wrapAround {
		if ret >= max {
			ret = max - secondsub
		}
		if ret < min {
			ret = min
		}
	} else {
		if ret >= max {
			ret = min + (ret - max)
		}
		if ret < min {
			ret = max - secondsub
		}
	}
	return ret
}

func Initialize() Model {
	m := Model{
		// TwelveHr: false,
		Seconds: true,
		value:   time.Now(),
		Styles: Styles{
			OuterBorder: lipgloss.NewStyle().
				Padding(1, 2).
				Align(lipgloss.Center).
				Margin(0, 2).
				Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#ff4884")),
			SelectedValue: lipgloss.NewStyle().Padding(0, 3).
				Foreground(lipgloss.Color("#000000")).Bold(true).
				Background(lipgloss.Color("#d55f87")),
			Value: lipgloss.NewStyle().Padding(0, 3),
		},
		selected: 0,
		keys:     keys,
	}
	m.Styles.SelectedOuterBorder = m.Styles.OuterBorder.Copy().
		Border(defBorder).BorderForeground(lipgloss.Color("#f48fb1"))
	return m
}
