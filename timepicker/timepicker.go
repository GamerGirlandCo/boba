package timepicker

import (
	"fmt"
	"log"
	"os"
	"time"

	util "git.tablet.sh/tablet/boba/types"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	z "github.com/lrstanley/bubblezone"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/term"
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

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.TickUp, k.TickDown, k.PrevField, k.NextField, k.Choose}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.TickUp, k.TickDown, k.PrevField, k.NextField},
		{k.First, k.Last, k.Choose},
	}
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
	First: key.NewBinding(
		key.WithKeys(tea.KeyHome.String(), tea.KeyShiftUp.String()),
		key.WithHelp("home/shift+↑", "first value"),
	),
	Last: key.NewBinding(
		key.WithKeys(tea.KeyEnd.String(), tea.KeyShiftDown.String()),
		key.WithHelp("end/shift+↓", "last value"),
	),
	Choose: key.NewBinding(
		key.WithKeys(tea.KeySpace.String(), tea.KeyEnter.String()),
		key.WithHelp("space/enter", "finalize selection"),
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
	help     help.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	thirdField := 1
	if m.Seconds {
		thirdField = 2
	}

	wo, _, _ := term.GetSize(int(os.Stdout.Fd()))
	m.width = wo
	// tea.Println("wja ", m.width)
	switch msg := msg.(type) {
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseLeft:
			for i := range headers {

				switch {
				case z.Get(fmt.Sprintf("%s-up", headers[i])).InBounds(msg):
					m.incOrDec(-1, 1, i)
				case z.Get(fmt.Sprintf("%s-dn", headers[i])).InBounds(msg):
					m.incOrDec(1, 1, i)
				}
			}
		case tea.MouseWheelUp:
			m.incOrDec(-1, 1, m.selected)
		case tea.MouseWheelDown:
			m.incOrDec(1, 1, m.selected)
		}
	case tea.KeyMsg:
		switch {
		// case "q", "ctrl+c":
		// 	return m, tea.Quit
		case key.Matches(msg, m.keys.NextField):
			m.selected = minMax(thirdField, m.selected, 1)
		case key.Matches(msg, m.keys.PrevField):
			m.selected = minMax(thirdField, m.selected, -1)
		case key.Matches(msg, m.keys.TickUp):
			m.incOrDec(-1, 1, m.selected)
		case key.Matches(msg, m.keys.TickDown):
			m.incOrDec(1, 1, m.selected)
		case key.Matches(msg, m.keys.Choose):
			return m, func() tea.Msg {
				// if twelvhr...
				result := m.value.Format("15:04:05")
				return util.GenResultMsg[string]{
					Res: result,
				}
			}
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
					z.Mark(fmt.Sprintf("%s-up", headers[a]), upArrowWide),
					bs.Render(jd),
					z.Mark(fmt.Sprintf("%s-dn", headers[a]), downArrowWide),
				)
			}
		}
	}
	log.Print("hi", m.width)
	return lipgloss.PlaceHorizontal(int(m.width), lipgloss.Center, lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.JoinHorizontal(lipgloss.Center, final...),
		m.help.View(m.keys),
	),
	)

}

func (m *Model) incOrDec(dir, amt, anotherarg int) {
	var multiplier time.Duration
	switch anotherarg {
	case 0:
		multiplier = time.Hour
	case 1:
		multiplier = time.Minute
	case 2:
		multiplier = time.Second
	}
	m.value = m.value.Add(multiplier * time.Duration(dir*amt))
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
	j := 0
	grad := genGradient()
	for i := -3; i < 4; i++ {
		var inten int
		if i < 0 {
			inten = 3 + i
		} else {
			inten = (7 - i - 3) - 1
		}

		// fmt.Println(inten)
		ihatethis := which + i + 1
		inter := minMaxTicks(0, wrapArg,
			ihatethis,
			-1, true, true)
		ret[j] = ticky{
			val:   inter,
			color: grad[inten],
		}
		j++
	}
	return ret
}

func genGradient() []string {
	ograd := make([]colorful.Color, 4)
	x0, _ := colorful.Hex("#444444")
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

func minMaxTicks(min, max, cur, dir int, wrapAround bool, inclusive bool) int {
	ret := cur
	if dir < 0 {
		ret--
	} else if dir > 0 {
		ret++
	}
	// ret += dir

	if ret >= max {
		ret = min + (ret - max)
	}
	if ret < min {
		ret = max + ret
	}
	return ret
}

func minMax(max, cur, dir int) int {
	cur += dir
	if cur < 0 {
		cur = max
	} else if cur > max {
		cur = 0
	}
	return cur
}

func Initialize(sekunti bool) Model {
	z.NewGlobal()
	helpo := help.New()
	helpo.ShowAll = true
	m := Model{
		// TwelveHr: false,
		Seconds: sekunti,
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
		help:     helpo,
	}
	m.Styles.SelectedOuterBorder = m.Styles.OuterBorder.Copy().
		Border(defBorder).BorderForeground(lipgloss.Color("#f48fb1"))
	return m
}
