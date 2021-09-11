package fetcher

import (
	"bytes"
	"io/ioutil"

	"golang.org/x/net/html"
)

type mockFetcher struct{}

func (m mockFetcher) Fetch(filePath string) (*html.Node, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(bytes.NewReader(file))
	if err != nil {
		return nil, err
	}

	return doc, nil
}
