package types

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
