package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const NEW_TAB_TITLE = "New Tab"

type tab struct {
	title  string
	url    string
	active bool
}

func getRegularBorder() lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = "┴"
	border.BottomRight = "┴"
	return border
}
func getActiveBorder() lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft, border.BottomRight = border.BottomRight, border.BottomLeft
	border.Bottom = ""
	return border
}

const maxTabWidth = 14

var (
	tabStyle lipgloss.Style = lipgloss.NewStyle().
			Foreground(TEXT_PRIMARY).Width(maxTabWidth).
			Padding(0, 1).Align(lipgloss.Center).
			Border(getRegularBorder()).BorderForeground(BORDER).
			Transform(func(s string) string {
			if len(s) >= maxTabWidth {
				s = s[:maxTabWidth-3] + "…"
			}
			return s
		})

	activeTabStyle lipgloss.Style = tabStyle.Border(getActiveBorder())
)

func (t tab) view() string {
	style := tabStyle
	if t.active {
		style = activeTabStyle
	}
	return style.Render(t.title)
}

type tabBar struct {
	ind   int
	tabs  []tab
	width int
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
	case tea.WindowSizeMsg:
		t.width = msg.Width

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
	var s []string
	for _, tab := range t.tabs {
		s = append(s, tab.view())
	}

	r := lipgloss.JoinHorizontal(lipgloss.Bottom, s...)

	remWidth := t.width - lipgloss.Width(r)
	if remWidth > 0 {
		r += lipgloss.NewStyle().Foreground(BORDER).Render(strings.Repeat("─", remWidth))
	}

	return r + "\n"
}

func (t tabBar) focusCmd(newTab bool) tea.Cmd {
	return func() tea.Msg {
		return tabChangedMsg{
			url:    t.tabs[t.ind].url,
			newTab: newTab,
		}
	}
}
