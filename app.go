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

// application is used to create and manage
// the application logic from the main thread.
type application struct {
	router  *gin.Engine
	port    int
	fetcher fetcher.Fetcher
}

// Init is used to initialize an application instance and its router.
// The mode parameter can only be gin.ReleaseMode or gin.DebugMode.
// If an invalid mode is provided, then the default gin.DebugMode
// will be used.
func (a *application) Init(port int, mode string, f fetcher.Fetcher) {
	// Set gin mode.
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	} else if mode != gin.DebugMode {
		log.Printf("[WARNING] Invalid application mode provided '%v', "+
			"starting up in default mode '%v'\n", mode, gin.DebugMode)
	}

	// Assign application variables.
	a.port = port
	a.router = gin.Default()
	a.fetcher = f

	// Load and set HTML templates.
	a.router.Static("/public", "public")
	a.router.LoadHTMLGlob("templates/*.html")

	// Define routes.
	a.router.GET("/", a.redirectToIndex)
	a.router.GET("/index", a.index)
	a.router.POST("/index", a.index)
	a.router.GET("/ping", a.ping)
}

// Start will attempt to run the router and listen for requests
// on the port provided during initialization.
func (a *application) Start() error {
	log.Println("[INFO] Starting application on port:", a.port)
	return a.router.Run(":" + strconv.Itoa(a.port))
}

// Stop can be used to execute any teardown functionality
// for the application when the service is to shutdown.
// Note that currently there is no teardown logic as it
// is currently not required.
func (a *application) Stop() error {
	log.Println("[INFO] Shutting down application")
	return nil
}

// redirectToIndex is used to redirect the base route
// to the index page.
func (a *application) redirectToIndex(c *gin.Context) {
	c.Redirect(http.StatusPermanentRedirect, "/index")
}

// index is the main route of the application,
// used to serve both GET and POST requests.
func (a *application) index(c *gin.Context) {
	parameters := gin.H{}
	parameters["title"] = "Web Analyser"
	parameters["stylesheet"] = "public/css/index.css"

	// If request method is of type GET, then return here.
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "index.html", parameters)

		return
	}

	// Obtain the URL entered by user from the form.
	url := c.PostForm("url")
	if url == "" {
		parameters["warning"] = "Website URL must not be empty"
		c.HTML(http.StatusBadRequest, "index.html", parameters)

		return
	}

	parameters["url"] = url

	// Check if obtained URL is valid.
	if u, err := netURL.Parse(url); err != nil || u.Host == "" || u.Scheme == "" {
		parameters["warning"] = "The provided URL is not valid, please enter a valid host and scheme"
		c.HTML(http.StatusBadRequest, "index.html", parameters)

		return
	}

	log.Println("[INFO] Provided URL:", url)

	// Fetch HTML document from the provided URL.
	document, err := a.fetcher.Fetch(url)
	if err != nil {
		parameters["warning"] = "The provided URL does not exist/is not accessible at the moment: " + err.Error()
		c.HTML(http.StatusInternalServerError, "index.html", parameters)

		return
	}

	// Crawl the fetched document, and return the results identified.
	info := crawl(document, url)

	log.Println("[INFO] Crawled Output:", info)

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

	log.Println("[INFO] Parameters Returned:", parameters)

	c.HTML(http.StatusOK, "index.html", parameters)
}

// ping is used for making sure that the application is running.
func (a *application) ping(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  http.StatusOK,
			"message": "Service is running",
		},
	)
}
