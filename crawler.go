package main

import (
	"golang.org/x/net/html"
)

func extractVerion(html string) (string, error) {
	return "", nil
}

func extractPageTitle(n *html.Node) (string, bool) {
	if n.Type == html.ElementNode && n.Data == "title" {
		return n.FirstChild.Data, true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		title, ok := extractPageTitle(c)
		if ok {
			return title, ok
		}
	}

	return "", false
}

func extractHeadingInfo(html string) (string, error) {
	return "", nil
}

func extractFormInfo(html string) (string, error) {
	return "", nil
}
