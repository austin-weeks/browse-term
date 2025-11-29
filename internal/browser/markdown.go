package browser

import (
	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/austin-weeks/browse-term/internal/themes"
	"github.com/charmbracelet/glamour"
	"golang.org/x/net/html"
)

func toMarkdown(n *html.Node) (string, error) {
	b, err := htmltomarkdown.ConvertNode(n)
	if err != nil {
		return "", err
	}
	md := string(b)
	if md == "" {
		md = "# Looks like this page is empty :(\n\nThere's nothing to see here - the site probably requires JavaScript to display its contents"
	}
	return md, nil
}

func RenderMarkdown(content string, w int, theme themes.Theme) (string, error) {
	r, err := glamour.NewTermRenderer(
		theme.RendererTheme(),
		glamour.WithWordWrap(w),
	)
	if err != nil {
		return "", err
	}
	return r.Render(content)
}
