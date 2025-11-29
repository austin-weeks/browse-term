package themes

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

const (
	BORDER        = lipgloss.Color("#7dd3fc")
	TextSecondary = lipgloss.Color("#9ca3af")
	TextTertiary  = lipgloss.Color("#4b5563")
)

func (t Theme) HighlightColor() lipgloss.Color {
	var color string
	switch t {
	case Dark:
		color = "39"
	case Light:
		color = "27"
	case TokyoNight:
		color = "#7aa2f7"
	case Pink:
		color = "212"
	case Dracula:
		color = "#bd93f9"

	case CatppuccinLatte:
		color = "#7287fd"
	case CatppuccinFrappe:
		color = "#8caaee"
	case CatppuccinMacchiato:
		color = "#8aadf4"
	case CatppuccinMocha:
		color = "#89b4fa"

	default:
		log.Error("Unknown theme value", "theme", t)
		color = "255"
	}
	return lipgloss.Color(color)
}

func (t Theme) TextPrimary() lipgloss.Color {
	var color string
	switch t {
	case Dark:
		color = "252"
	case Light:
		color = "234"
	case TokyoNight:
		color = "#a9b1d6"
	case Pink:
		color = "255"
	case Dracula:
		color = "#f8f8f2"

	case CatppuccinLatte:
		color = "#4c4f69"
	case CatppuccinFrappe:
		color = "#c6d0f5"
	case CatppuccinMacchiato:
		color = "#cad3f5"
	case CatppuccinMocha:
		color = "#cdd6f4"

	default:
		log.Error("Unknown theme value", "theme", t)
		color = "255"
	}
	return lipgloss.Color(color)
}
