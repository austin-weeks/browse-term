package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type searchFocusLostMsg struct{}
type searchConfirmedMsg struct {
	url string
}

type searchBar struct {
	style lipgloss.Style
	input textinput.Model
}

func newSearchBar() searchBar {
	input := textinput.New()
	// input.Placeholder = "Type a URL"
	input.Focus()
	input.Prompt = ""

	return searchBar{
		style: lipgloss.NewStyle(),
		input: input,
	}
}

func (s searchBar) Init() tea.Cmd {
	return textinput.Blink
}

func (s searchBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case shouldFocusMsg:
		s.input.Focus()
		return s, textinput.Blink
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEscape:
			s.input.Blur()
			return s, func() tea.Msg { return searchFocusLostMsg{} }
		case tea.KeyEnter:
			s.input.Blur()
			return s, func() tea.Msg {
				return searchConfirmedMsg{
					url: s.input.Value(),
				}
			}
		}
	}
	s.input, cmd = s.input.Update(msg)
	cmds = append(cmds, cmd)

	return s, tea.Batch(cmds...)
}

func (s searchBar) View() string {
	return s.input.View() + "\n"
}
