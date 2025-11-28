package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const prompt = "https:// "

type searchBar struct {
	style lipgloss.Style
	input textinput.Model
}

func newSearchBar() searchBar {
	input := textinput.New()
	input.Prompt = prompt
	input.PromptStyle = input.PromptStyle.Foreground(TEXTSECONDARY).PaddingLeft(1)

	style := lipgloss.NewStyle().BorderForeground(BORDER).Border(lipgloss.RoundedBorder()).Foreground(TEXTPRIMARY)

	return searchBar{
		style: style,
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
	case tea.WindowSizeMsg:
		withBorders := msg.Width - 3
		s.style = s.style.Width(withBorders)
		s.input.Width = withBorders - len(prompt) - 2

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
				func() tea.Msg { return triggerFetchMsg{url: s.input.Value()} },
				s.loseFocus(),
			)
		default:
			s.input, cmd = s.input.Update(msg)
			cmds = append(cmds, cmd)
		}

	case pageContentMsg:
		s.input.SetValue(msg.c.URL)

	default:
		s.input, cmd = s.input.Update(msg)
		cmds = append(cmds, cmd)
	}

	return s, tea.Batch(cmds...)
}

func (s searchBar) View() string {
	v := s.input.View()
	return s.style.Render(v) + "\n"
}

func (s searchBar) loseFocus() tea.Cmd {
	return func() tea.Msg {
		return focusChangedMsg{
			target: focusPage,
		}
	}
}
