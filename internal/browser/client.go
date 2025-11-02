package browser

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// TODO - if a page errors, and we nav to another tab, then we go back, then the page is blank :(

// TODO - clear cache after certain time to avoid memory growing endlessly (probably not something to worry about)

// url -> WebPage
var webPages map[string]WebPage = make(map[string]WebPage)

// Represents a website page
type WebPage struct {
	Url     string
	Title   string
	Content string
	Links   []string
}

func FetchWebPage(url string) (WebPage, error) {
	var prettyUrl string
	if strings.HasPrefix(url, "http://") {
		prettyUrl = strings.TrimLeft(url, "http://")
	} else if strings.HasPrefix(url, "https://") {
		prettyUrl = strings.TrimLeft(url, "https://")
	} else {
		prettyUrl = url
		url = "https://" + url
	}

	if page, ok := webPages[url]; ok {
		return page, nil
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
	if dom, err := html.Parse(strings.NewReader(page)); err == nil {
		if t, err := extractTitle(dom); err == nil {
			title = t
		}
	}

	webPage := WebPage{
		Url:     prettyUrl,
		Title:   title,
		Content: page,
		Links:   []string{},
	}
	webPages[url] = webPage
	return webPage, nil
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
