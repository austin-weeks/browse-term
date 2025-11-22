// Package browser provides methods for fetching websites for display in browse-term.
package browser

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"golang.org/x/net/html"
)

// url -> WebPage
var webPages map[string]*WebPage = make(map[string]*WebPage)

// WebPage represents a website page
type WebPage struct {
	URL     string
	Title   string
	Content string
	Links   []Link
}

func FetchWebPage(path string, enableJS bool) (*WebPage, error) {
	path, _ = strings.CutPrefix(path, "https://")
	path, _ = strings.CutPrefix(path, "http://")
	prettyURL := path
	path = "https://" + path

	URL, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	url := URL.String()

	if page, ok := webPages[url]; ok {
		return page, nil
	}

	var dom *html.Node
	if enableJS {
		dom, err = fetchWithChrome(url)
		if err != nil {
			return nil, err
		}
	} else {
		dom, err = fetchPlain(url)
		if err != nil {
			return nil, err
		}
	}

	var title string
	if t, err := extractTitle(dom); err == nil {
		title = t
	} else {
		title = url
	}

	md, err := toMarkdown(dom)
	if err != nil {
		return nil, err
	}

	links := extractLinks(dom, URL)

	webPage := &WebPage{
		URL:     prettyURL,
		Title:   title,
		Content: md,
		Links:   links,
	}
	webPages[url] = webPage
	return webPage, nil
}

func fetchPlain(url string) (*html.Node, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Browse-Term (https://github.com/austin-weeks/browse-term)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // nolint

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Non-200 response: %v", resp.StatusCode)
	}

	return html.Parse(resp.Body)
}

func fetchWithChrome(url string) (*html.Node, error) {
	// Create context with headless Chrome
	ctx, ogCancel := chromedp.NewContext(context.Background())
	defer ogCancel()

	// Set timeout
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var htmlStr string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.OuterHTML("html", &htmlStr),
	)
	if err != nil {
		return nil, err
	}

	return html.Parse(strings.NewReader(htmlStr))
}
