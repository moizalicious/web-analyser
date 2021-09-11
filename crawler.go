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
	"net/url"

	"golang.org/x/net/html"
)

type pageInfo struct {
	htmlVersion               string
	pageTitle                 string
	headingCount              headingCounts
	accessibleInternalLinks   linkInfo
	unaccessibleInternalLinks linkInfo
	externalLinks             linkInfo
	containsForm              bool
}

type headingCounts struct {
	h1 int
	h2 int
	h3 int
	h4 int
	h5 int
	h6 int
}

type linkInfo struct {
	count int
	links []string
}

type link struct {
	href         string
	isAccessible bool
	isExternal   bool
}

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

const (
	htmlV401Strict       = "-//W3C//DTD HTML 4.01//EN"
	htmlV401Transitional = "-//W3C//DTD HTML 4.01 Transitional//EN"
	htmlV401Frameset     = "-//W3C//DTD HTML 4.01 Frameset//EN"
	xhtmlV1Strict        = "-//W3C//DTD XHTML 1.0 Strict//EN"
	xhtmlV1Transitional  = "-//W3C//DTD XHTML 1.0 Transitional//EN"
	xhtmlV1Frameset      = "-//W3C//DTD XHTML 1.0 Frameset//EN"
	xhtmlV11             = "-//W3C//DTD XHTML 1.1//EN"
)

func crawl(document *html.Node, host string) pageInfo {
	info := pageInfo{}
	info.accessibleInternalLinks.links = make([]string, 0)
	info.unaccessibleInternalLinks.links = make([]string, 0)
	info.externalLinks.links = make([]string, 0)

	var crawler func(*html.Node)

	crawler = func(n *html.Node) {
		switch n.Type {
		case html.ElementNode:
			switch n.Data {
			case anchor:
				l, err := identifyLinkInfo(n.Attr, host)
				if err != nil {
					log.Printf("[WARNING] Failed to obtain link information from anchor element: %v\n", err)
				} else if l.isExternal {
					info.externalLinks.links = append(info.externalLinks.links, l.href)
					info.externalLinks.count++
				} else if l.isAccessible {
					info.accessibleInternalLinks.links = append(info.accessibleInternalLinks.links, l.href)
					info.accessibleInternalLinks.count++
				} else {
					info.unaccessibleInternalLinks.links = append(info.unaccessibleInternalLinks.links, l.href)
					info.unaccessibleInternalLinks.count++
				}

			case form:
				info.containsForm = true

			case title:
				if n.FirstChild != nil {
					info.pageTitle = n.FirstChild.Data
				}

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

func identifyLinkInfo(attributes []html.Attribute, host string) (link, error) {
	for _, a := range attributes {
		if a.Key == "href" {
			isAccessible, isExternal, err := extractLinkInfo(a.Val, host)
			if err != nil {
				return link{}, err
			}

			l := link{}
			l.href = a.Val
			l.isAccessible = isAccessible
			l.isExternal = isExternal

			return l, nil
		}
	}

	return link{}, errors.New("no href attribute available in given list")
}

// first bool isAccessible
// second bool isExternal
func extractLinkInfo(href string, host string) (bool, bool, error) {
	h, err := url.Parse(href)
	if err != nil || h.Host == "" || h.Scheme == "" {
		// href link is not accessible, therefore it is not even external
		return false, false, nil
	}

	x, err := url.Parse(host)
	if err != nil || h.Host == "" || h.Scheme == "" {
		// comparer link must be valid, ideally this should never happen
		return false, false, errors.New("provided host url is invalid")
	}

	if h.Host == x.Host {
		// if both hosts are the same, then the link is accessible but not external
		return true, false, nil
	} else {
		// if both are different, then the link is accessible and external
		return true, true, nil
	}
}
