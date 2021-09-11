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
	"log"
	"net/http"
	netURL "net/url"
	"strconv"
	"web-analyser/fetcher"

	"github.com/gin-gonic/gin"
)

type application struct {
	router  *gin.Engine
	port    int
	fetcher fetcher.Fetcher
}

func (a *application) Init(port int, mode string, f fetcher.Fetcher) {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	} else if mode != gin.DebugMode {
		log.Printf("[WARNING] Invalid $APP_MODE environment variable defined '%v', "+
			"starting up in default mode '%v'", mode, gin.DebugMode)
	}

	a.port = port
	a.router = gin.Default()
	a.fetcher = f

	a.router.Static("/public", "public")
	a.router.LoadHTMLGlob("templates/*.html")

	a.router.GET("/", a.redirectToIndex)
	a.router.GET("/index", a.index)
	a.router.POST("/index", a.index)
	a.router.GET("/ping", a.ping)
}

func (a *application) Start() error {
	return a.router.Run(":" + strconv.Itoa(a.port))
}

func (a *application) Stop() {
	// TODO - add some teardown functionality
}

func (a *application) redirectToIndex(c *gin.Context) {
	c.Redirect(http.StatusPermanentRedirect, "/index")
}

func (a *application) index(c *gin.Context) {
	parameters := gin.H{}
	parameters["title"] = "Web Analyser"
	parameters["stylesheet"] = "public/css/index.css"

	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "index.html", parameters)

		return
	}

	url := c.PostForm("url")
	if url == "" {
		parameters["warning"] = "Website URL must not be empty"
		c.HTML(http.StatusBadRequest, "index.html", parameters)

		return
	}

	parameters["url"] = url

	if u, err := netURL.Parse(url); err != nil || u.Host == "" || u.Scheme == "" {
		parameters["warning"] = "The provided URL is not valid"
		c.HTML(http.StatusBadRequest, "index.html", parameters)

		return
	}

	log.Println("Provided URL:", url)

	document, err := a.fetcher.Fetch(url)
	if err != nil {
		parameters["warning"] = "The provided URL does not exist/is not accessible at the moment"
		c.HTML(http.StatusInternalServerError, "index.html", parameters)

		return
	}

	info := crawl(document, url)
	log.Println("Crawled Output:", info)

	parameters["displayResult"] = true
	parameters["htmlVersion"] = info.htmlVersion
	parameters["pageTitle"] = info.pageTitle

	parameters["h1Count"] = info.headingCount.h1
	parameters["h2Count"] = info.headingCount.h2
	parameters["h3Count"] = info.headingCount.h3
	parameters["h4Count"] = info.headingCount.h4
	parameters["h5Count"] = info.headingCount.h5
	parameters["h6Count"] = info.headingCount.h6

	parameters["accessibleInternalLinkCount"] = info.accessibleInternalLinks.count
	parameters["accessibleInternalLinks"] = info.accessibleInternalLinks.links

	parameters["unaccessibleInternalLinkCount"] = info.unaccessibleInternalLinks.count
	parameters["unaccessibleInternalLinks"] = info.unaccessibleInternalLinks.links

	parameters["externalLinkCount"] = info.externalLinks.count
	parameters["externalLinks"] = info.externalLinks.links

	parameters["containsForm"] = info.containsForm

	c.HTML(http.StatusOK, "index.html", parameters)
}

func (a *application) ping(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  http.StatusOK,
			"message": "Service is running",
		},
	)
}
