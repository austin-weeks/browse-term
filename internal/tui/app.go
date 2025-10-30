package tui

import (
	"github.com/austin-weeks/browse-term/internal/browser"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type focus int

const (
	searchFocus focus = iota
	pageFocus
)

type app struct {
	// State
	cols, rows int
	focus      focus

	// Components
	title lipgloss.Style
	// Tabs later
	searchBar tea.Model
	page      tea.Model
	keybinds  keybinds
}

func New() app {
	return app{
		focus: searchFocus,

		title:    lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).Foreground(lipgloss.Color("#080808")).SetString("Terminal Browser"),
		keybinds: newKeybinds(),

		searchBar: newSearchBar(),
		page:      newPage(),
	}
}

func (a app) Init() tea.Cmd {
	return tea.Batch(
		a.page.Init(),
		a.searchBar.Init(),
	)
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
		a.keybinds.setWidth(a.cols)

		msg.Height -= lipgloss.Height(a.title.Render() + "\n")
		msg.Height -= lipgloss.Height(a.keybinds.view(a.focus))

		a.searchBar, cmd = a.searchBar.Update(msg)
		cmds = append(cmds, cmd)
		msg.Height -= lipgloss.Height(a.searchBar.View() + "\n")

		a.page, cmd = a.page.Update(msg)
		cmds = append(cmds, cmd)

	case tea.KeyMsg:
		if s := msg.String(); s == "ctrl+c" || (a.focus != searchFocus && s == "q") {
			return a, tea.Quit
		}
		if a.focus == searchFocus {
			a.searchBar, cmd = a.searchBar.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			if msg.String() == "/" {
				a.focus = searchFocus
				a.searchBar, cmd = a.searchBar.Update(shouldFocusMsg{})
				cmds = append(cmds, cmd)
			}

			a.page, cmd = a.page.Update(msg)
			cmds = append(cmds, cmd)
		}

	case focusLostMsg:
		switch msg.focus {
		case searchFocus:
			a.focus = pageFocus
		case pageFocus:
			// no op
		}

	case searchConfirmedMsg:
		resp, err := browser.FetchWebPage(msg.url)
		if err != nil {
			cmds = append(cmds, func() tea.Msg { return pageErrMsg{err: err} })
		}
		cmds = append(cmds, func() tea.Msg { return pageContentMsg{contents: resp} })

	case pageContentMsg:
		// println(msg.contents.Title)
		// We need to send the title to the tabs
		a.page, cmd = a.page.Update(msg)
		cmds = append(cmds, cmd)
		a.searchBar, cmd = a.searchBar.Update(msg)
		cmds = append(cmds, cmd)

	case pageErrMsg:
		a.page, cmd = a.page.Update(msg)
		cmds = append(cmds, cmd)

	default:
	}

	return a, tea.Batch(cmds...)
}

func (a app) View() string {
	s := a.title.Render() + "\n"
	s += a.searchBar.View()
	s += a.page.View()
	s += a.keybinds.view(a.focus)

	return s
}
