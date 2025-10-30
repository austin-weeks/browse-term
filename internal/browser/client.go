package browser

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// Represents a website page
type WebPage struct {
	Url     string
	Title   string
	Content string
	Links   []string
}

func FetchWebPage(url string) (WebPage, error) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	resp, err := http.Get(url)
	if err != nil {
		return WebPage{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return WebPage{}, fmt.Errorf("Non-200 response: %v", resp.StatusCode)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return WebPage{}, err
	}
	page := string(bytes)

	title := url
	dom, err := html.Parse(strings.NewReader(page))
	if err == nil {
		t, err := extractTitle(dom)
		if err == nil {
			title = t
		}
	}

	return WebPage{
		Url:     url,
		Title:   title,
		Content: page,
		Links:   []string{},
	}, nil
}

func extractTitle(parent *html.Node) (string, error) {
	// Base Case
	if parent.Type == html.ElementNode && parent.Data == "title" {
		var title string
		for child := parent.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.TextNode {
				title += child.Data
			}
		}
		return strings.TrimSpace(title), nil
	}

	for child := parent.FirstChild; child != nil; child = child.NextSibling {
		title, err := extractTitle(child)
		if err == nil && title != "" {
			return title, nil
		}
	}

	return "", fmt.Errorf("title element not found")
}
