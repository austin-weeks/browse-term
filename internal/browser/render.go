package browser

import (
	"fmt"
	"regexp"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/austin-weeks/browse-term/internal/themes"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/log"
	"golang.org/x/net/html"
)

var mdImgRegex = regexp.MustCompile(`!\[.*?\]\(.*?\)`)

const mdImgPlaceholder = "{{MD_IMAGE}}"

func toMarkdown(n *html.Node) (string, error) {
	b, err := htmltomarkdown.ConvertNode(n)
	if err != nil {
		return "", err
	}
	md := string(b)
	if md == "" {
		md = "# Looks like this page is empty :(\n\nThere's nothing to see here - the site probably requires JavaScript to display its contents"
	}
	md = mdImgRegex.ReplaceAllString(md, mdImgPlaceholder)
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

func (wp *WebPage) RenderPage(w int, theme themes.Theme) (string, error) {
	r, err := glamour.NewTermRenderer(
		theme.RendererTheme(),
		glamour.WithWordWrap(w),
	)
	if err != nil {
		return "", err
	}

	sections := strings.Split(wp.Content, mdImgPlaceholder)
	if len(sections) != len(wp.images)+1 {
		panic(fmt.Sprintf("%d sections and %d images", len(sections), len(wp.images)))
	}

	// Combine sections and images
	var builder strings.Builder
	for i := range len(sections) - 1 {
		s, err := r.Render(sections[i])
		if err != nil {
			return "", err
		}
		builder.WriteString(s)
		img := wp.images[i]
		s, err = wp.images[i].render(w)
		if err != nil {
			s = fmt.Sprintf("![%s](%s)", img.alt, img.src)
			log.Error("Failed to render image", "error", err, "src", img.src)
		}
		builder.WriteString(s)
	}
	s, err := r.Render(sections[len(sections)-1])
	if err != nil {
		return "", err
	}
	builder.WriteString(s)

	return builder.String(), nil
}
