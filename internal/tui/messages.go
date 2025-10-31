package tui

import "github.com/austin-weeks/browse-term/internal/browser"

type focus int

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

type tabChangedMsg struct {
	url    string
	newTab bool
}
