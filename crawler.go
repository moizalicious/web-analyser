/*
   Copyright 2021 github.com/moizalicious

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package main

import (
	"errors"
	"log"
	netURL "net/url"

	"golang.org/x/net/html"
)

// documentInfo is used to structure the information of a
// HTML document that can be identified by the crawler.
type documentInfo struct {
	htmlVersion               string
	pageTitle                 string
	headingCount              headingCounts
	accessibleInternalLinks   linkInfo
	inaccessibleInternalLinks linkInfo
	externalLinks             linkInfo
	containsForm              bool
}

// headingCounts is used to keep track of the count of
// all the Heading levels in a HTML document.
type headingCounts struct {
	h1 int
	h2 int
	h3 int
	h4 int
	h5 int
	h6 int
}

// linkInfo is used to keep track of a set of
// anchor links in a HTML document.
type linkInfo struct {
	count int
	links []string
}

// link is used to define the properties of
// a single anchor link href value.
type link struct {
	href         string
	isAccessible bool
	isExternal   bool
}

// HTML element names.
const (
	anchor = "a"
	form   = "form"
	title  = "title"

	heading1 = "h1"
	heading2 = "h2"
	heading3 = "h3"
	heading4 = "h4"
	heading5 = "h5"
	heading6 = "h6"
)

// HTML Doctype Versions.
const (
	htmlV401Strict       = "-//W3C//DTD HTML 4.01//EN"
	htmlV401Transitional = "-//W3C//DTD HTML 4.01 Transitional//EN"
	htmlV401Frameset     = "-//W3C//DTD HTML 4.01 Frameset//EN"
	xhtmlV1Strict        = "-//W3C//DTD XHTML 1.0 Strict//EN"
	xhtmlV1Transitional  = "-//W3C//DTD XHTML 1.0 Transitional//EN"
	xhtmlV1Frameset      = "-//W3C//DTD XHTML 1.0 Frameset//EN"
	xhtmlV11             = "-//W3C//DTD XHTML 1.1//EN"
)

// crawl is used to crawl through a provided HTML document and identify
// a given set of properties. These properties are defined and returned
// as a documentInfo struct.
//
// Aside from the document, a url string needs to be provided which defines
// the url in which the document was obtained from. The reason for this
// is so that the crawler can use this to differentiate between internally
// and externally accessible links.
func crawl(document *html.Node, url string) documentInfo {
	// Initialise documentInfo struct.
	info := documentInfo{}
	info.accessibleInternalLinks.links = make([]string, 0)
	info.inaccessibleInternalLinks.links = make([]string, 0)
	info.externalLinks.links = make([]string, 0)

	// Define crawler function.
	var crawler func(*html.Node)

	// Initialize crawler function.
	crawler = func(n *html.Node) {
		// Identify node type.
		switch n.Type {
		// Element node.
		case html.ElementNode:
			// Identify element type.
			switch n.Data {
			// <a> element.
			case anchor:
				l, err := identifyLinkInfo(n.Attr, url)
				if err != nil {
					log.Println("[WARNING] Failed to obtain link information from anchor element:", err)
				} else if l.isExternal {
					info.externalLinks.links = append(info.externalLinks.links, l.href)
					info.externalLinks.count++
				} else if l.isAccessible {
					info.accessibleInternalLinks.links = append(info.accessibleInternalLinks.links, l.href)
					info.accessibleInternalLinks.count++
				} else {
					info.inaccessibleInternalLinks.links = append(info.inaccessibleInternalLinks.links, l.href)
					info.inaccessibleInternalLinks.count++
				}
			// <form> element.
			case form:
				info.containsForm = true
			// <title> element.
			case title:
				if n.FirstChild != nil {
					info.pageTitle = n.FirstChild.Data
				}
			// <h1> to <h6> elements.
			case heading1:
				info.headingCount.h1++
			case heading2:
				info.headingCount.h2++
			case heading3:
				info.headingCount.h3++
			case heading4:
				info.headingCount.h4++
			case heading5:
				info.headingCount.h5++
			case heading6:
				info.headingCount.h6++
			}
		// Doctype node.
		case html.DoctypeNode:
			info.htmlVersion = identifyHTMLVersion(n.Attr)
		}

		// Crawl through children/next sibling of current node.
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawler(c)
		}
	}

	// Execute crawler function.
	crawler(document)

	return info
}

// identifyHTMLVersion returns the HTML version of a document that is
// identified from the attributes of a html.DoctypeNode only.
func identifyHTMLVersion(attributes []html.Attribute) string {
	if len(attributes) == 0 {
		return "HTML 5"
	} else {
		for _, attr := range attributes {
			if attr.Key == "public" {
				switch attr.Val {
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

// identifyLinkInfo identifies the properties of a href attribute of a anchor element.
// If the provided list of attributes does not contain a href value, then an error is returned.
func identifyLinkInfo(attributes []html.Attribute, url string) (link, error) {
	for _, attr := range attributes {
		if attr.Key == "href" {
			isAccessible, isExternal, err := extractLinkInfo(attr.Val, url)
			if err != nil {
				return link{}, err
			}

			l := link{}
			l.href = attr.Val
			l.isAccessible = isAccessible
			l.isExternal = isExternal

			return l, nil
		}
	}

	return link{}, errors.New("no href attribute available in given list")
}

// extractLinkInfo is used to identify if a href value is
// accessible or inaccessible, and internal or external.
func extractLinkInfo(href string, url string) (isAccessible bool, isExternal bool, err error) {
	h, err := netURL.Parse(href)
	if err != nil || h.Host == "" || h.Scheme == "" {
		// href link is not accessible, therefore it is not
		// considered as external either.
		return false, false, nil
	}

	u, err := netURL.Parse(url)
	if err != nil || h.Host == "" || h.Scheme == "" {
		// comparer url must be valid, ideally this should never happen.
		return false, false, errors.New("provided host url is invalid")
	}

	if h.Host == u.Host {
		// if both hosts are the same, then the link is accessible but not external.
		return true, false, nil
	} else {
		// if both hosts are different, then the link is accessible and external.
		return true, true, nil
	}
}
