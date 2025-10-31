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

func newSearchBar() searchBar {
	input := textinput.New()
	input.Placeholder = "Type a URL"
	input.Prompt = ""

	return searchBar{
		style: lipgloss.NewStyle(),
		input: input,
	}
}

func (s searchBar) Init() tea.Cmd {
	return nil
}

func (s searchBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case focusChangedMsg:
		if msg.target != focusSearch {
			break
		}
		cmd = s.input.Focus()
		cmds = append(cmds, cmd)
		s.input.CursorEnd()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEscape:
			s.input.Blur()
			cmds = append(cmds, s.loseFocus())
		case tea.KeyEnter:
			if s.input.Value() == "" {
				break
			}
			s.input.Blur()
			cmds = append(cmds,
				func() tea.Msg { return searchConfirmedMsg{url: s.input.Value()} },
				s.loseFocus(),
			)
		default:
			s.input, cmd = s.input.Update(msg)
			cmds = append(cmds, cmd)
		}

	case tabChangedMsg:
		s.input.SetValue(msg.url)
		if msg.newTab {
			cmds = append(cmds, func() tea.Msg {
				return focusChangedMsg{
					target: focusSearch,
				}
			})
		}

	case pageContentMsg:
		s.input.SetValue(msg.c.Url)

	default:
		s.input, cmd = s.input.Update(msg)
		cmds = append(cmds, cmd)
	}

	return s, tea.Batch(cmds...)
}

func (s searchBar) View() string {
	return s.input.View() + "\n"
}

func (s searchBar) loseFocus() tea.Cmd {
	return func() tea.Msg {
		return focusChangedMsg{
			target: focusPage,
		}
	}
}
