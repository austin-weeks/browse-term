package tui

import (
	"github.com/austin-weeks/browse-term/internal/browser"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var tableBaseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(BORDER)

type linkTable struct {
	table table.Model
	links []browser.Link
}

func newLinkTable() linkTable {
	columns := []table.Column{
		{Title: "Link", Width: 40},
		{Title: "URL", Width: 60},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(false),
		table.WithHeight(18),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(BORDER).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color(BORDER)).
		Background(GREY600).
		Bold(false)
	t.SetStyles(s)

	return linkTable{
		table: t,
		links: []browser.Link{},
	}
}

func (l linkTable) Init() tea.Cmd {
	return nil
}

func (l linkTable) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return l, asCmd(focusChangedMsg{target: focusPage})

		case "enter":
			if len(l.links) > 0 {
				url := l.links[l.table.Cursor()].URL
				cmds = append(cmds, asCmd(focusChangedMsg{target: focusPage}))
				cmds = append(cmds, asCmd(triggerFetchMsg{url: url}))
				l.table.GotoTop()
			}
		default:
			l.table, cmd = l.table.Update(msg)
			cmds = append(cmds, cmd)
		}

	case pageContentMsg:
		l.updateRows(msg.c.Links)

	case focusChangedMsg:
		if msg.target == focusLinks {
			l.table.Focus()
		} else {
			l.table.Blur()
		}

	default:
		l.table, cmd = l.table.Update(msg)
		cmds = append(cmds, cmd)
	}

	return l, tea.Batch(cmds...)
}

func (l linkTable) View() string {
	return tableBaseStyle.Render(l.table.View()) + "\n"
}

func (l *linkTable) updateRows(links []browser.Link) {
	l.links = links
	var rows []table.Row
	for _, link := range l.links {
		name := link.Name
		rows = append(rows, table.Row{
			name,
			link.URL,
		})
	}
	l.table.SetRows(rows)
}
