package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var pageFocusKeys = []string{
	"q - quit",
	"/ - search",
}

var searchFocusKeys = []string{
	"^c - quit",
	"esc - exit search",
	"enter - go to URl",
}

type keybinds struct {
	style lipgloss.Style
}

func newKeybinds() keybinds {
	return keybinds{
		style: lipgloss.NewStyle().Foreground(lipgloss.Color(TEXT_LIGHT)).Align(lipgloss.Center),
	}
}

func (k *keybinds) setWidth(w int) {
	k.style = k.style.Width(w)
}

func (k keybinds) view(focus focus) string {
	var keys []string
	switch focus {
	case pageFocus:
		keys = pageFocusKeys
	case searchFocus:
		keys = searchFocusKeys
	default:
		panic("unhandled focus")
	}

	return k.style.Render(strings.Join(keys, "       "))
}
