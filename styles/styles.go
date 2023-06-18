package styles

import lg "github.com/charmbracelet/lipgloss"

type Theme struct {
	Text             lg.Color
	ActiveText       lg.Color
	ActiveBackground lg.Color
	FaintText        lg.Color
	FaintBackground  lg.Color
	TitleText        lg.Color
	TitleBackground  lg.Color
	Border           lg.Color
	ActiveBorder     lg.Color
	InactiveBorder   lg.Color
}

type Styles struct {
	Text           lg.Style
	Active         lg.Style
	Faint          lg.Style
	Title          lg.Style
	Border         lg.Style
	ActiveBorder   lg.Style
	InactiveBorder lg.Style
}

var DefaultTheme = Theme{
	Text:             lg.Color("#ffffff"),
	ActiveText:       lg.Color("#000000"),
	ActiveBackground: lg.Color("#00e7e3"),
	FaintText:        lg.Color("#bbbbbb"),
	FaintBackground:  lg.Color("0"),
	TitleText:        lg.Color("#ff4884"),
	TitleBackground:  lg.Color("0"),
	ActiveBorder:     lg.Color("#00e7e3"),
	InactiveBorder:   lg.Color("#aaaaaa"),
	Border:           lg.Color("#a3fffd"),
}

var DefaultStyles = Styles{
	Text:   lg.NewStyle().Foreground(DefaultTheme.Text),
	Active: lg.NewStyle().Foreground(DefaultTheme.ActiveText).Background(DefaultTheme.ActiveBackground),
}
