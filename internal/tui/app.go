package tui

import (
	"github.com/austin-weeks/browse-term/internal/browser"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type app struct {
	// State
	focus focus

	// Components
	title     lipgloss.Style
	tabBar    tea.Model
	searchBar tea.Model
	page      tea.Model
	keybinds  keybinds
}

// Get a new browse-term application
func New() app {
	title := lipgloss.NewStyle().Bold(true).Align(lipgloss.Center).
		Foreground(TEXT_PRIMARY).SetString("Terminal Browser").
		BorderBottom(true).BorderStyle(lipgloss.NormalBorder()).BorderForeground(BORDER)

	return app{
		focus: focusSearch,

		title:     title,
		tabBar:    newTabBar(),
		searchBar: newSearchBar(),
		page:      newPage(),
		keybinds:  newKeybinds(),
	}
}

func (a app) Init() tea.Cmd {
	return tea.Batch(
		a.tabBar.Init(),
		a.searchBar.Init(),
		a.page.Init(),
		func() tea.Msg {
			return focusChangedMsg{
				target: focusSearch,
			}
		},
	)
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		msg.Height += 3
		a.title = a.title.Width(msg.Width)
		a.keybinds.setWidth(msg.Width)

		msg.Height -= lipgloss.Height(a.title.Render() + "\n")
		msg.Height -= lipgloss.Height(a.keybinds.view(a.focus))

		a.tabBar, cmd = a.tabBar.Update(msg)
		cmds = append(cmds, cmd)
		msg.Height -= lipgloss.Height(a.tabBar.View())

		a.searchBar, cmd = a.searchBar.Update(msg)
		cmds = append(cmds, cmd)
		msg.Height -= lipgloss.Height(a.searchBar.View())

		a.page, cmd = a.page.Update(msg)
		cmds = append(cmds, cmd)

	case tea.KeyMsg:
		if s := msg.String(); s == "ctrl+c" || (a.focus != focusSearch && s == "q") {
			return a, tea.Quit
		}

		if a.focus == focusSearch {
			a.searchBar, cmd = a.searchBar.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			if msg.String() == "/" {
				cmds = append(cmds, func() tea.Msg {
					return focusChangedMsg{
						target: focusSearch,
					}
				})
			}
			a.tabBar, cmd = a.tabBar.Update(msg)
			cmds = append(cmds, cmd)

			a.page, cmd = a.page.Update(msg)
			cmds = append(cmds, cmd)
		}

	case focusChangedMsg:
		a.focus = msg.target
		cmds = append(cmds, a.updateAllComponents(msg))

	case searchConfirmedMsg:
		resp, err := browser.FetchWebPage(msg.url)
		if err != nil {
			cmd = func() tea.Msg { return pageErrMsg{err: err} }
		} else {
			cmd = func() tea.Msg { return pageContentMsg{c: resp} }
		}
		cmds = append(cmds, cmd)

	case pageContentMsg:
		a.tabBar, cmd = a.tabBar.Update(msg)
		cmds = append(cmds, cmd)
		a.searchBar, cmd = a.searchBar.Update(msg)
		cmds = append(cmds, cmd)
		a.page, cmd = a.page.Update(msg)
		cmds = append(cmds, cmd)

	case pageErrMsg:
		a.tabBar, cmd = a.tabBar.Update(msg)
		cmds = append(cmds, cmd)
		a.page, cmd = a.page.Update(msg)
		cmds = append(cmds, cmd)

	case tabChangedMsg:
		if msg.url != "" {
			cmds = append(cmds, func() tea.Msg {
				return searchConfirmedMsg{
					url: msg.url,
				}
			})
			break
		}
		a.searchBar, cmd = a.searchBar.Update(msg)
		cmds = append(cmds, cmd)

		a.page, cmd = a.page.Update(pageContentMsg{
			c: browser.WebPage{Content: ""}},
		)
		cmds = append(cmds, cmd)

	default:
		cmd = a.updateAllComponents(msg)
		cmds = append(cmds, cmd)
	}

	return a, tea.Batch(cmds...)
}

func (a app) View() string {
	s := a.title.Render() + "\n"
	s += a.tabBar.View()
	s += a.searchBar.View()
	s += a.page.View()
	s += a.keybinds.view(a.focus)

	return s
}

func (a *app) updateAllComponents(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	a.tabBar, cmd = a.tabBar.Update(msg)
	cmds = append(cmds, cmd)

	a.searchBar, cmd = a.searchBar.Update(msg)
	cmds = append(cmds, cmd)

	a.page, cmd = a.page.Update(msg)
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}
