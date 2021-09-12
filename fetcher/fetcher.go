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

package fetcher

import (
	"golang.org/x/net/html"
)

// Fetcher provides the interface with a function
// which can be used to fetch and parse HTML content
// from a given source.
type Fetcher interface {
	Fetch(string) (*html.Node, error)
}

// NewURLFetcher creates and returns a instance of urlFetcher.
func NewURLFetcher() Fetcher {
	return urlFetcher{}
}

// NewURLFetcher creates and returns a instance of fileFetcher.
func NewFileFetcher() Fetcher {
	return fileFetcher{}
}
