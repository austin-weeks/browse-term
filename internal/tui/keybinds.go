package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

type keyMap []key.Binding

func (k keyMap) ShortHelp() []key.Binding {
	return k
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k}
}

var pageFocusKeys keyMap = []key.Binding{
	key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
	key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
	key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "open links")),
	key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("↓/j", "scroll down")),
	key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("↑/k", "scroll up")),
}

var searchFocusKeys keyMap = []key.Binding{
	key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("^c", "quit")),
	key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "exit search")),
	key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "fetch site")),
}

var linksFocusKeys keyMap = []key.Binding{
	key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "exit links")),
	key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "open link")),
	key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("↓/j", "down")),
	key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("↑/k", "up")),
}

type keybinds struct {
	width int
	help  help.Model
}

func newKeybinds() keybinds {
	help := help.New()
	help.Styles.ShortKey = help.Styles.ShortKey.Foreground(TEXTPRIMARY)
	help.Styles.ShortDesc = help.Styles.ShortDesc.Foreground(TEXTSECONDARY)
	help.Styles.ShortSeparator = help.Styles.ShortSeparator.Foreground(TEXTSECONDARY)
	return keybinds{
		help: help,
	}
}

func (k *keybinds) setWidth(w int) {
	k.width = w
	k.help.Width = w
}

func (k keybinds) view(focus focus) string {
	var keys keyMap
	switch focus {
	case focusPage:
		keys = pageFocusKeys
	case focusSearch:
		keys = searchFocusKeys
	case focusLinks:
		keys = linksFocusKeys
	default:
		panic("unhandled focus")
	}

	s := k.help.ShortHelpView(keys) + "\n"
	return lipgloss.PlaceHorizontal(k.width, lipgloss.Center, s)
}
