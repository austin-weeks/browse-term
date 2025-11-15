package browser

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Name string
	URL  string
}

func extractLinks(node *html.Node, baseURL *url.URL) []Link {
	var links []Link
	if node.Type == html.ElementNode && node.Data == "a" {
		var path string
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				path = attr.Val
			}
		}
		if path == "" {
			return nil
		}
		u, err := url.Parse(path)
		if err == nil && baseURL != nil {
			u = baseURL.ResolveReference(u)
		}
		link := Link{
			URL: u.String(),
		}

		var name string
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.TextNode {
				name += child.Data
			}
		}
		link.Name = strings.TrimSpace(name)
		return append(links, link)
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		links = append(links, extractLinks(child, baseURL)...)
	}

	return links
}

func extractTitle(node *html.Node) (string, error) {
	// Base Case
	if node.Type == html.ElementNode && node.Data == "title" {
		var title string
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.TextNode {
				title += child.Data
			}
		}
		return strings.TrimSpace(title), nil
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		title, err := extractTitle(child)
		if err == nil && title != "" {
			return title, nil
		}
	}

	return "", fmt.Errorf("title element not found")
}
