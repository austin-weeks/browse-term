package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

func welcomeScreen(w int, h int) string {
	return lipgloss.NewStyle().Width(w).Height(h).
		Foreground(BORDER).AlignHorizontal(lipgloss.Center).AlignVertical(lipgloss.Center).
		Render(`
 /$$$$$$$                                                     /$$$$$$$$                               
| $$__  $$                                                   |__  $$__/                               
| $$  \ $$  /$$$$$$   /$$$$$$  /$$  /$$  /$$  /$$$$$$$  /$$$$$$ | $$  /$$$$$$   /$$$$$$  /$$$$$$/$$$$ 
| $$$$$$$  /$$__  $$ /$$__  $$| $$ | $$ | $$ /$$_____/ /$$__  $$| $$ /$$__  $$ /$$__  $$| $$_  $$_  $$
| $$__  $$| $$  \__/| $$  \ $$| $$ | $$ | $$|  $$$$$$ | $$$$$$$$| $$| $$$$$$$$| $$  \__/| $$ \ $$ \ $$
| $$  \ $$| $$      | $$  | $$| $$ | $$ | $$ \____  $$| $$_____/| $$| $$_____/| $$      | $$ | $$ | $$
| $$$$$$$/| $$      |  $$$$$$/|  $$$$$/$$$$/ /$$$$$$$/|  $$$$$$$| $$|  $$$$$$$| $$      | $$ | $$ | $$
|_______/ |__/       \______/  \_____/\___/ |_______/  \_______/|__/ \_______/|__/      |__/ |__/ |__/
`)
}

func loadingScreen(w int, h int) string {
	return lipgloss.NewStyle().Width(w).Height(h).
		Foreground(BORDER).AlignHorizontal(lipgloss.Center).AlignVertical(lipgloss.Center).
		Render(`o o o`)
}

type page struct {
	ready     bool
	viewport  viewport.Model
	contentFn func(w int, h int) (string, error)
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
		if !p.ready {
			p.ready = true
			p.setContent(func(w int, h int) (string, error) {
				return welcomeScreen(w-5, h), nil
			})
		} else {
			cmds = append(cmds, p.setContent(p.contentFn))
		}

	case tea.KeyMsg:
		p.viewport, cmd = p.viewport.Update(msg)
		cmds = append(cmds, cmd)

	case pageContentMsg:
		if msg.c.Content == "" {
			p.setContent(func(w int, h int) (string, error) {
				return welcomeScreen(w, h), nil
			})
		} else {
			cmd = p.setContent(func(w int, h int) (string, error) {
				return renderMarkdown(msg.c.Content, w)
			})
			cmds = append(cmds, cmd)
		}

	case pageErrMsg:
		p.viewport.SetContent("Something went wrong :(\n\n\n" + msg.err.Error())

	case onLoadMsg:
		p.viewport.SetContent(loadingScreen(p.viewport.Width-4, p.viewport.Height))

	default:
		p.viewport, cmd = p.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return p, tea.Batch(cmds...)
}

func (p page) View() string {
	return p.viewport.View() + "\n"
}

func renderMarkdown(content string, w int) (string, error) {
	r, err := glamour.NewTermRenderer(
		glamour.WithStylePath("dark"),
		glamour.WithWordWrap(w),
	)
	if err != nil {
		return "", err
	}
	return r.Render(content)
}

func (p *page) setContent(fn func(w int, h int) (string, error)) tea.Cmd {
	p.contentFn = fn
	s, err := fn(p.viewport.Width, p.viewport.Height)
	if err != nil {
		return asCmd(pageErrMsg{err: err})
	}
	p.viewport.SetContent(s)
	return nil
}
