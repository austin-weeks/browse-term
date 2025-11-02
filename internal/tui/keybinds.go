package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var pageFocusKeys = []string{
	"q - quit",
	"/ - search",
	"h - tab left",
	"l - tab right",
	"j - scroll down",
	"k - scroll up",
	"^w - close tab",
	"^t - new tab",
}

var searchFocusKeys = []string{
	"^c - quit",
	"esc - exit search",
	"enter - go to URL",
}

type keybinds struct {
	style lipgloss.Style
}

func newKeybinds() keybinds {
	return keybinds{
		style: lipgloss.NewStyle().Foreground(TEXT_SECONDARY).Align(lipgloss.Center),
	}
}

func (k *keybinds) setWidth(w int) {
	k.style = k.style.Width(w)
}

func (k keybinds) view(focus focus) string {
	var keys []string
	switch focus {
	case focusPage:
		keys = pageFocusKeys
	case focusSearch:
		keys = searchFocusKeys
	default:
		panic("unhandled focus")
	}

	return k.style.Render(strings.Join(keys, "       "))
}
