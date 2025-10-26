package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type page struct {
	viewport viewport.Model
	style    lipgloss.Style
}

func newPage() page {
	viewport := viewport.New(0, 0)
	viewport.Style = viewport.Style.Border(lipgloss.DoubleBorder())
	return page{
		style:    lipgloss.NewStyle().Background(lipgloss.Color("black")).Border(lipgloss.NormalBorder()),
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
		p.style = p.style.Width(msg.Width).Height(msg.Height)
		p.viewport.Width, p.viewport.Height = msg.Width, msg.Height

	case tea.KeyMsg:
		p.viewport, cmd = p.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}
	return p, tea.Batch(cmds...)
}

func (p page) View() string {
	return p.viewport.View()
}
