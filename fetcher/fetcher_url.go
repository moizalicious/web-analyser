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
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

type urlFetcher struct{}

func (u urlFetcher) Fetch(url string) (*html.Node, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("http status %v returned", response.Status)
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Printf("[WARNING] Failed to close response body: %v\n", err)
		}
	}()

	document, err := html.Parse(response.Body)
	if err != nil {
		return nil, err
	}

	return document, nil
}
