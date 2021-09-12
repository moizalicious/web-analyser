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
	"testing"
	"web-analyser/fetcher"

	"github.com/stretchr/testify/assert"
)

func TestCrawl(t *testing.T) {
	type testCase struct {
		name         string
		fileLocation string
		mockURL      string
		expected     documentInfo
	}

	f := fetcher.NewFileFetcher()

	testCases := []testCase{
		{
			name:         "basic_crawl_test",
			fileLocation: "res/test.html",
			mockURL:      "https://www.google.com",
			expected: documentInfo{
				htmlVersion: "HTML 5",
				pageTitle:   "Test Page",
				headingCount: headingCounts{
					h1: 1,
					h2: 1,
					h3: 1,
					h4: 1,
					h5: 1,
					h6: 1,
				},
				accessibleInternalLinks: linkInfo{
					count: 2,
					links: []string{
						"https://www.google.com",
						"https://www.google.com/test",
					},
				},
				inaccessibleInternalLinks: linkInfo{
					count: 3,
					links: []string{"#Home", "/test", "/"},
				},
				externalLinks: linkInfo{
					count: 2,
					links: []string{
						"https://www.w3schools.com",
						"https://www.youtube.com",
					},
				},
				containsForm: true,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			document, err := f.Fetch(testCase.fileLocation)
			if err != nil {
				t.Error("Failed to fetch document from given file location:", err)
			}

			actual := crawl(document, testCase.mockURL)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func BenchmarkCrawl(b *testing.B) {
	sampleURL := "https://www.google.com"
	document, err := fetcher.NewFileFetcher().Fetch("res/test.html")
	if err != nil {
		b.Error("Failed to fetch HTML content of the provided file:", err)
	}

	for i := 0; i < b.N; i++ {
		crawl(document, sampleURL)
	}
}
