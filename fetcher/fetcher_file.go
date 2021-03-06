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
	"bytes"
	"io/ioutil"

	"golang.org/x/net/html"
)

// fileFetcher is the implementation of Fetcher to
// fetch and parse an HTML document from a local file.
type fileFetcher struct{}

// Fetch and parse a HTML document in the given path.
func (f fileFetcher) Fetch(filePath string) (*html.Node, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	document, err := html.Parse(bytes.NewReader(file))
	if err != nil {
		return nil, err
	}

	return document, nil
}
