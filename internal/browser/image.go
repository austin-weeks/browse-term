package browser

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/austin-weeks/browse-term/internal/chafa"
)

// TODO: track width and height attributes
type image struct {
	src   string
	title string
	alt   string
	data  []byte
	w, h  int
}

func fetchImages(images []image) error {
	ch := make(chan error, len(images))
	var wg sync.WaitGroup
	for i := range len(images) {
		img := &images[i]
		wg.Go(func() { img.fetch(ch) })
	}
	wg.Wait()
	close(ch)

	var errs []error
	for err := range ch {
		errs = append(errs, err)
	}
	if len(errs) == 0 {
		return nil
	} else {
		return errors.Join(errs...)
	}
}

func (i *image) fetch(ch chan error) {
	// TODO: need to handle data urls separately
	req, err := http.NewRequest("GET", i.src, nil)
	if err != nil {
		ch <- err
		return
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "image/*")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ch <- fmt.Errorf("failed to fetch image from %s: %w", i.src, err)
		return
	}
	defer resp.Body.Close() // nolint

	if resp.StatusCode != http.StatusOK {
		ch <- fmt.Errorf("non-200 response: %v", resp.StatusCode)
		return
	}
	p, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- err
		return
	}
	i.data = p
	ch <- nil
}

func (i *image) render(width int) (string, error) {
	if len(i.data) == 0 {
		return "", fmt.Errorf("image has no data")
	}

	// TODO: use image width from html attribute

	return chafa.RenderImage(i.data, width)
}
