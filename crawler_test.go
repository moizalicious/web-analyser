package main

import (
	"testing"
	"web-analyser/fetcher"
)

func TestCrawl(t *testing.T) {
	f := fetcher.NewMockFetcher()

	document, err := f.Fetch("res/test.html")
	if err != nil {
		t.Error(err)
	}

	info := crawl(document)

	t.Log("OUTPUT:", info)
}
