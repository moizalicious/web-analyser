package fetcher

import (
	"bytes"
	"io/ioutil"

	"golang.org/x/net/html"
)

type mockFetcher struct{}

func (m mockFetcher) Fetch(url string) (*html.Node, error) {
	file, err := ioutil.ReadFile("res/test.html")
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(bytes.NewReader(file))
	if err != nil {
		return nil, err
	}

	return doc, nil
}
