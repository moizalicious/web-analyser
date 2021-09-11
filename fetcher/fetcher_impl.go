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

type fetcher struct{}

func (f fetcher) Fetch(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error when calling http get: %v", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("[WARNING] Failed to close response body: %v\n", err)
		}
	}()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
