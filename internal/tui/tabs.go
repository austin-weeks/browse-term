package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const NEW_TAB_TITLE = "Enter a URL"

type tab struct {
	title  string
	url    string
	active bool
}

type tabBar struct {
	ind   int
	tabs  []tab
	style lipgloss.Style
}

func newTabBar() tabBar {
	return tabBar{
		ind: 0,
		tabs: []tab{
			{
				title:  NEW_TAB_TITLE,
				active: true,
			},
		},
	}
}

func (t tabBar) Init() tea.Cmd {
	return nil
}

func (t tabBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "h":
			if len(t.tabs) == 1 {
				break
			}
			t.tabs[t.ind].active = false
			t.ind--
			if t.ind < 0 {
				t.ind = len(t.tabs) - 1
			}
			t.tabs[t.ind].active = true
			cmds = append(cmds, t.focusCmd(false))

		case "l":
			if len(t.tabs) == 1 {
				break
			}
			t.tabs[t.ind].active = false
			t.ind++
			if t.ind >= len(t.tabs) {
				t.ind = 0
			}
			t.tabs[t.ind].active = true
			cmds = append(cmds, t.focusCmd(false))

		case "ctrl+t":
			t.tabs[t.ind].active = false
			newTab := tab{
				title:  NEW_TAB_TITLE,
				active: true,
			}
			t.tabs = append(t.tabs, newTab)
			t.ind = len(t.tabs) - 1
			cmds = append(cmds, t.focusCmd(true))

		case "ctrl+w":
			newTab := false
			if len(t.tabs) == 1 {
				t.tabs[t.ind] = tab{
					title:  NEW_TAB_TITLE,
					active: true,
				}
				newTab = true
			} else {
				t.tabs = append(t.tabs[:t.ind], t.tabs[t.ind+1:]...)
				t.ind--
				t.tabs[t.ind].active = true
			}
			cmds = append(cmds, t.focusCmd(newTab))
		}

	case pageErrMsg:
		tab := &t.tabs[t.ind]
		tab.title = "Error"

	case pageContentMsg:
		tab := &t.tabs[t.ind]
		tab.url = msg.c.Url
		tab.title = msg.c.Title
	}

	return t, tea.Batch(cmds...)
}

func (t tabBar) View() string {
	s := strings.Builder{}
	for _, tab := range t.tabs {
		if tab.active {
			s.WriteString("* " + tab.title + "   ")
		} else {
			s.WriteString(tab.title + "   ")
		}
	}
	s.WriteString("\n")
	return t.style.Render(s.String())
}

func (t tabBar) focusCmd(newTab bool) tea.Cmd {
	return func() tea.Msg {
		return tabChangedMsg{
			url:    t.tabs[t.ind].url,
			newTab: newTab,
		}
	}
}
