package main

import (
	"golang.org/x/net/html"
)

type pageInfo struct {
	htmlVersion       string
	pageTitle         string
	headings          map[element]int
	links             []link
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

type link struct {
}

func crawl(document *html.Node) pageInfo {
	info := pageInfo{}
	info.headings = make(map[element]int)

	var crawler func(*html.Node)

	crawler = func(n *html.Node) {
		switch n.Type {
		case html.ElementNode:
			switch n.Data {
			case string(anchor):
				// TODO - get link information

			case string(form):
				info.containsLoginForm = true

			case string(title):
				info.pageTitle = n.FirstChild.Data

			case string(heading1):
				info.headings[heading1]++
			case string(heading2):
				info.headings[heading2]++
			case string(heading3):
				info.headings[heading3]++
			case string(heading4):
				info.headings[heading4]++
			case string(heading5):
				info.headings[heading5]++
			case string(heading6):
				info.headings[heading6]++
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
		for _, v := range attributes {
			if v.Key == "public" {
				switch v.Val {
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
