package tui

import (
	"time"

	"github.com/austin-weeks/browse-term/internal/browser"
	"github.com/austin-weeks/browse-term/internal/themes"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func welcomeScreen(w int, h int, theme themes.Theme) string {
	return lipgloss.NewStyle().Width(w).Height(h).
		Foreground(theme.HighlightColor()).AlignHorizontal(lipgloss.Center).AlignVertical(lipgloss.Center).
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

type page struct {
	theme themes.Theme

	ready          bool
	refreshContent bool
	viewport       viewport.Model
	spinner        spinner.Model
	contentFn      func(w int, h int, p page) (string, error)
}

func newPage(jsEnabled bool, theme themes.Theme) page {
	viewport := viewport.New(0, 0)
	viewport.Style = viewport.Style.Border(lipgloss.RoundedBorder()).BorderForeground(theme.HighlightColor())
	spin := spinner.New()
	if jsEnabled {
		spin.Spinner = spinner.Spinner{
			Frames: []string{
				" Rendering  ",
				" Rendering  ",
				" Rendering  ",
				" Rendering  ",
				" Rendering  ",
				" Rendering  ",
			},
			FPS: time.Second / 13,
		}
	} else {
		spin.Spinner = spinner.Spinner{
			Frames: []string{"∙∙∙∙∙∙∙", "●∙∙∙∙∙∙", "∙●∙∙∙∙∙", "∙∙●∙∙∙∙", "∙∙∙●∙∙∙", "∙∙∙∙●∙∙", "∙∙∙∙∙●∙", "∙∙∙∙∙∙●"},
			FPS:    time.Second / 9,
		}
	}

	return page{
		theme:    theme,
		viewport: viewport,
		spinner:  spin,
	}
}

func (p page) Init() tea.Cmd {
	return p.spinner.Tick
}

func (p page) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.viewport.Width, p.viewport.Height = msg.Width, msg.Height
		if !p.ready {
			p.ready = true
			p.setContent(func(w int, h int, p page) (string, error) {
				return welcomeScreen(w-5, h, p.theme), nil
			})
		} else {
			cmds = append(cmds, p.setContent(p.contentFn))
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "g":
			p.viewport.GotoTop()
		case "G":
			p.viewport.GotoBottom()
		}
		p.viewport, cmd = p.viewport.Update(msg)
		cmds = append(cmds, cmd)

	case pageContentMsg:
		p.viewport.GotoTop()
		p.refreshContent = false
		if msg.c.Content == "" {
			p.setContent(func(w int, h int, p page) (string, error) {
				return welcomeScreen(w, h, p.theme), nil
			})
		} else {
			cmd = p.setContent(func(w int, h int, p page) (string, error) {
				return browser.RenderMarkdown(msg.c.Content, w, p.theme)
			})
			cmds = append(cmds, cmd)
		}

	case pageErrMsg:
		p.viewport.GotoTop()
		p.refreshContent = false
		p.setContent(func(w int, h int, p page) (string, error) {
			text := "# Something went wrong :(\n\n\n" + msg.err.Error()
			s, err := browser.RenderMarkdown(text, w, p.theme)
			if err != nil {
				return text, nil
			}
			return s, nil
		})

	case onLoadMsg:
		p.viewport.GotoTop()
		p.refreshContent = true
		p.setContent(func(w, h int, p page) (string, error) {
			s := p.spinner.View()
			s = lipgloss.Place(w-8, h, lipgloss.Center, lipgloss.Center, s)
			return s, nil
		})

	default:
		p.viewport, cmd = p.viewport.Update(msg)
		cmds = append(cmds, cmd)
		p.spinner, cmd = p.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return p, tea.Batch(cmds...)
}

func (p page) View() string {
	if p.refreshContent {
		p.setContent(p.contentFn)
	}
	return p.viewport.View() + "\n"
}

func (p *page) setContent(fn func(w int, h int, p page) (string, error)) tea.Cmd {
	p.contentFn = fn
	s, err := fn(p.viewport.Width, p.viewport.Height, *p)
	if err != nil {
		return asCmd(pageErrMsg{err: err})
	}
	p.viewport.SetContent(s)
	return nil
}
