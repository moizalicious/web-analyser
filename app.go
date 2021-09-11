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
	Router  *gin.Engine
	Port    int
	Fetcher fetcher.Fetcher
}

func (a *application) Init(port int, f fetcher.Fetcher) {
	a.Port = port
	a.Router = gin.Default()
	a.Fetcher = f

	a.Router.Static("/public", "public")
	a.Router.LoadHTMLGlob("templates/*.html")

	a.Router.GET("/", a.redirectToIndex)
	a.Router.GET("/index", a.index)
	a.Router.POST("/index", a.index)
	a.Router.GET("/ping", a.ping)
}

func (a *application) Start() error {
	return a.Router.Run(":" + strconv.Itoa(a.Port))
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

	if _, err := netURL.ParseRequestURI(url); err != nil {
		parameters["warning"] = "The provided URL is not valid"
		c.HTML(http.StatusBadRequest, "index.html", parameters)

		return
	}

	log.Println("Provided URL:", url)

	document, err := a.Fetcher.Fetch(url)
	if err != nil {
		parameters["warning"] = "The provided URL does not exist/is not accessible at the moment"
		c.HTML(http.StatusInternalServerError, "index.html", parameters)

		return
	}

	info := crawl(document)
	log.Println("Crawled Output:", info)

	parameters["displayResult"] = true
	parameters["htmlVersion"] = info.htmlVersion
	parameters["pageTitle"] = info.pageTitle

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
