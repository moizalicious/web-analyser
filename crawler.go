package main

import (
	"errors"
	"log"
	"net/url"

	"golang.org/x/net/html"
)

type pageInfo struct {
	htmlVersion       string
	pageTitle         string
	headingCount      map[element]int
	internalLinkCount int
	internalLinks     []string
	externalLinkCount int
	externalLinks     []string
	containsLoginForm bool
}

type element string

const (
	anchor element = "a"
	form   element = "form"
	title  element = "title"

	heading1 element = "h1"
	heading2 element = "h2"
	heading3 element = "h3"
	heading4 element = "h4"
	heading5 element = "h5"
	heading6 element = "h6"
)

const (
	htmlV401Strict       = "-//W3C//DTD HTML 4.01//EN"
	htmlV401Transitional = "-//W3C//DTD HTML 4.01 Transitional//EN"
	htmlV401Frameset     = "-//W3C//DTD HTML 4.01 Frameset//EN"
	xhtmlV1Strict        = "-//W3C//DTD XHTML 1.0 Strict//EN"
	xhtmlV1Transitional  = "-//W3C//DTD XHTML 1.0 Transitional//EN"
	xhtmlV1Frameset      = "-//W3C//DTD XHTML 1.0 Frameset//EN"
	xhtmlV11             = "-//W3C//DTD XHTML 1.1//EN"
)

func crawl(document *html.Node) pageInfo {
	info := pageInfo{}
	info.headingCount = make(map[element]int)
	info.internalLinks = make([]string, 0)
	info.externalLinks = make([]string, 0)

	var crawler func(*html.Node)

	crawler = func(n *html.Node) {
		switch n.Type {
		case html.ElementNode:
			switch n.Data {
			case string(anchor):
				href, isExteral, err := identifyLinkInfo(n.Attr)
				if err != nil {
					log.Printf("[WARNING] Failed to obtain link information from anchor element: %v\n", err)
				} else if isExteral {
					info.externalLinks = append(info.externalLinks, href)
					info.externalLinkCount++
				} else {
					info.internalLinks = append(info.internalLinks, href)
					info.internalLinkCount++
				}

			case string(form):
				info.containsLoginForm = true

			case string(title):
				if n.FirstChild != nil {
					info.pageTitle = n.FirstChild.Data
				}

			case string(heading1):
				info.headingCount[heading1]++
			case string(heading2):
				info.headingCount[heading2]++
			case string(heading3):
				info.headingCount[heading3]++
			case string(heading4):
				info.headingCount[heading4]++
			case string(heading5):
				info.headingCount[heading5]++
			case string(heading6):
				info.headingCount[heading6]++
			}

		case html.DoctypeNode:
			info.htmlVersion = identifyHTMLVersion(n.Attr)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawler(c)
		}
	}

	crawler(document)

	return info
}

func identifyHTMLVersion(attributes []html.Attribute) string {
	if len(attributes) == 0 {
		return "HTML 5"
	} else {
		for _, a := range attributes {
			if a.Key == "public" {
				switch a.Val {
				case htmlV401Strict:
					return "HTML v4.01 Strict"
				case htmlV401Transitional:
					return "HTML v4.01 Transitional"
				case htmlV401Frameset:
					return "HTML v4.01 Frameset"
				case xhtmlV1Strict:
					return "XHTML v1 Strict"
				case xhtmlV1Transitional:
					return "XHTML v1 Transitional"
				case xhtmlV1Frameset:
					return "XHTML v1 Frameset"
				case xhtmlV11:
					return "XHTML v1.1"
				default:
					return ""
				}
			}
		}
	}

	return ""
}

func identifyLinkInfo(attributes []html.Attribute) (string, bool, error) {
	for _, a := range attributes {
		if a.Key == "href" {
			return a.Val, isExternalLink(a.Val), nil
		}
	}

	return "", false, errors.New("no href attribute available in given list")
}

func isExternalLink(href string) bool {
	_, err := url.ParseRequestURI(href)
	return err != nil
}
