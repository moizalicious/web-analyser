package fetcher

import (
	"golang.org/x/net/html"
)

type Fetcher interface {
	Fetch(url string) (*html.Node, error)
}

func NewFetcher() Fetcher {
	return fetcher{}
}

func NewMockFetcher() Fetcher {
	return mockFetcher{}
}
