package main

import (
	"testing"
	"web-analyser/fetcher"

	"github.com/stretchr/testify/assert"
)

func TestExtractPageTitle(t *testing.T) {
	f := fetcher.NewMockFetcher()

	doc, err := f.Fetch("http://localhost:8080")
	if err != nil {
		t.Error(err)
	}

	expected := "Test File Title"
	actual, _ := extractPageTitle(doc)

	assert.Equal(t, expected, actual)
}
