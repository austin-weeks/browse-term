package tui

import "github.com/austin-weeks/browse-term/internal/browser"

type shouldFocusMsg struct{}

type focusLostMsg struct {
	focus focus
}

type searchConfirmedMsg struct {
	url string
}

type pageErrMsg struct {
	err error
}

type pageContentMsg struct {
	contents browser.WebPage
}
