package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type searchBar struct {
	style lipgloss.Style
	input textinput.Model
}

func searchFocusLostCmd() tea.Msg {
	return focusLostMsg{focus: searchFocus}
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
		case tea.KeyEscape, tea.KeyEnter:
			s.input.Blur()
			cmds = append(cmds, searchFocusLostCmd)

			if msg.Type == tea.KeyEnter {
				cmds = append(cmds, func() tea.Msg {
					return searchConfirmedMsg{
						url: s.input.Value(),
					}
				}, searchFocusLostCmd)
			}
		default:
			s.input, cmd = s.input.Update(msg)
			cmds = append(cmds, cmd)
		}

	case pageContentMsg:
		s.input.SetValue(msg.contents.Url)

	default:
		s.input, cmd = s.input.Update(msg)
		cmds = append(cmds, cmd)
	}

	return s, tea.Batch(cmds...)
}

func (s searchBar) View() string {
	return s.input.View() + "\n"
}
