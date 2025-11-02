package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type page struct {
	viewport viewport.Model
}

func newPage() page {
	viewport := viewport.New(0, 0)
	viewport.Style = viewport.Style.Border(lipgloss.RoundedBorder()).BorderForeground(BORDER)
	return page{
		viewport: viewport,
	}
}
func (p page) Init() tea.Cmd {
	return nil
}

func (p page) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.viewport.Width, p.viewport.Height = msg.Width, msg.Height

	case tea.KeyMsg:
		p.viewport, cmd = p.viewport.Update(msg)
		cmds = append(cmds, cmd)

	case pageContentMsg:
		p.viewport.SetContent(msg.c.Content)

	case pageErrMsg:
		p.viewport.SetContent("Something went wrong :(\n\n\n" + msg.err.Error())
	}
	return p, tea.Batch(cmds...)
}

func (p page) View() string {
	return p.viewport.View() + "\n"
}
