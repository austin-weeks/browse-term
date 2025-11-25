package browser

import (
	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
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
