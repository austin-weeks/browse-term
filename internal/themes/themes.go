// Package themes specifies themes for browse-term.
package themes

import (
	_ "embed"
	"fmt"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/log"
)

type Theme string

const (
	Dark       Theme = "dark"
	Light      Theme = "light"
	TokyoNight Theme = "tokyo-night"
	Pink       Theme = "pink"
	Dracula    Theme = "dracula"

	CatppuccinLatte     Theme = "catppuccin-latte"
	CatppuccinFrappe    Theme = "catppuccin-frappe"
	CatppuccinMacchiato Theme = "catppuccin-macchiato"
	CatppuccinMocha     Theme = "catppuccin-mocha"
)

//go:embed catppuccin-latte.json
var catppuccinLatte string

//go:embed catppuccin-frappe.json
var catppuccinFrappe string

//go:embed catppuccin-macchiato.json
var catppuccinMacchiato string

//go:embed catppuccin-mocha.json
var catppuccinMocha string

func (t Theme) RendererTheme() glamour.TermRendererOption {
	switch t {
	case Dark, Light, TokyoNight, Pink, Dracula:
		return glamour.WithStandardStyle(string(t))

	case CatppuccinLatte:
		return glamour.WithStylesFromJSONBytes([]byte(catppuccinLatte))
	case CatppuccinFrappe:
		return glamour.WithStylesFromJSONBytes([]byte(catppuccinFrappe))
	case CatppuccinMacchiato:
		return glamour.WithStylesFromJSONBytes([]byte(catppuccinMacchiato))
	case CatppuccinMocha:
		return glamour.WithStylesFromJSONBytes([]byte(catppuccinMocha))

	default:
		log.Error("Unknown theme, falling back to 'tokyo-night'", "theme", t)
		return glamour.WithStandardStyle(string(TokyoNight))
	}
}

func ValidateTheme(theme Theme) error {
	switch theme {
	case Dark, Light, TokyoNight, Pink, Dracula:
		return nil
	case CatppuccinLatte, CatppuccinFrappe, CatppuccinMacchiato, CatppuccinMocha:
		return nil
	default:
		return fmt.Errorf("Theme value '%s' is not a valid theme", theme)
	}
}
