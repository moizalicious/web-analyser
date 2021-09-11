package fetcher

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

type fetcher struct{}

func (f fetcher) Fetch(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error when calling http get: %v", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("[WARN] Failed to close response body: error - %v\n", err)
		}
	}()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
