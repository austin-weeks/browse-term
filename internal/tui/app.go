package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type app struct {
	cols, rows int
	title      gloss.Style
	// Tabs
	// SearchBar
	page tea.Model
	// Commands
}

func New() app {
	return app{
		page:  newPage(),
		title: gloss.NewStyle().Bold(true).Align(gloss.Center).Foreground(gloss.Color("#080808")).SetString("Terminal Browser"),
	}
}

func (a app) Init() tea.Cmd {
	a.page.Init()
	return nil
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.cols, a.rows = msg.Width, msg.Height
		a.title = a.title.Width(a.cols)
		msg.Height -= gloss.Height(a.title.Render() + "\n")
		a.page, cmd = a.page.Update(msg)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return a, tea.Quit
		}
		a.page, cmd = a.page.Update(msg)
		cmds = append(cmds, cmd)
	}

	return a, tea.Batch(cmds...)
}

func (a app) View() string {
	s := a.title.Render() + "\n"

	s += a.page.View() + "\n"

	return s
}
