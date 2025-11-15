package tui

import (
	"github.com/austin-weeks/browse-term/internal/browser"
	tea "github.com/charmbracelet/bubbletea"
)

type focus int

type onLoadMsg struct{}

const (
	focusSearch focus = iota
	focusPage
	focusLinks
)

type focusChangedMsg struct {
	target focus
}

type triggerFetchMsg struct {
	url string
}

type pageErrMsg struct {
	err error
}

// Sent with the focused page changes
type pageContentMsg struct {
	c browser.WebPage
}

func asCmd(msg any) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}
