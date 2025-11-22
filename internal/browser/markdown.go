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

	// s := "# Hello World!\n\nThis is some markdown.\n\nHere is some _italic text_, some **bold text**, and even some `code`!\n\nHere we've even got a code snippet:\n```typescript\nimport { safeTry } from \"errgo-ts\";\n\nconst res = safeTry(() => fetch(\"/api/users\"));\nif (res.err) {\n  console.error(err);\n}\n```\n\nHere we have an image - will this work?"
}
