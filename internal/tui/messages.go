package tui

import (
	"github.com/austin-weeks/browse-term/internal/browser"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

type focus int

type onLoadMsg struct{}

const (
	focusSearch focus = iota
	focusPage
)

type focusChangedMsg struct {
	target focus
}

type searchConfirmedMsg struct {
	url string
}

type pageErrMsg struct {
	err error
}

// Sent with the focused page changes
type pageContentMsg struct {
	c browser.WebPage
}

func (p pageContentMsg) renderMarkdown(width int) (string, error) {
	r, err := glamour.NewTermRenderer(
		glamour.WithStylePath("dark"),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return "", err
	}
	return r.Render(p.c.Content)
}

type tabChangedMsg struct {
	url    string
	newTab bool
}

func asCmd(msg any) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}
