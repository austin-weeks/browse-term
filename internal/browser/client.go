// Package browser provides methods for fetching websites for display in browse-term.
package browser

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// url -> WebPage
var webPages map[string]WebPage = make(map[string]WebPage)

// WebPage represents a website page
type WebPage struct {
	URL     string
	Title   string
	Content string
	Links   []Link
}

func FetchWebPage(path string) (WebPage, error) {
	path, _ = strings.CutPrefix(path, "https://")
	path, _ = strings.CutPrefix(path, "http://")
	prettyURL := path
	path = "https://" + path

	URL, err := url.Parse(path)
	if err != nil {
		return WebPage{}, err
	}
	url := URL.String()

	if page, ok := webPages[url]; ok {
		return page, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return WebPage{}, err
	}
	req.Header.Set("User-Agent", "Browse-Term (https://github.com/austin-weeks/browse-term)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return WebPage{}, err
	}
	defer resp.Body.Close() // nolint

	if resp.StatusCode != http.StatusOK {
		return WebPage{}, fmt.Errorf("Non-200 response: %v", resp.StatusCode)
	}

	var title string
	dom, err := html.Parse(resp.Body)
	if err == nil {
		if t, err := extractTitle(dom); err == nil {
			title = t
		}
	} else {
		title = url
	}

	md, err := toMarkdown(dom)
	if err != nil {
		return WebPage{}, err
	}

	links := extractLinks(dom, URL)

	webPage := WebPage{
		URL:     prettyURL,
		Title:   title,
		Content: md,
		Links:   links,
	}
	webPages[url] = webPage
	return webPage, nil
}
