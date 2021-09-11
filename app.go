package main

import (
	"log"
	"net/http"
	"strconv"
	"web-analyser/fetcher"

	"github.com/gin-gonic/gin"
)

type Application struct {
	Router  *gin.Engine
	Port    int
	Fetcher fetcher.Fetcher
}

func (a *Application) Init(port int, f fetcher.Fetcher) {
	a.Port = port
	a.Router = gin.Default()
	a.Fetcher = f

	a.Router.LoadHTMLGlob("templates/*.html")

	a.Router.GET("/", a.redirectToIndex)
	a.Router.GET("/index", a.index)
	a.Router.POST("/index", a.index)
	a.Router.GET("/ping", a.ping)
}

func (a *Application) Start() error {
	return a.Router.Run(":" + strconv.Itoa(a.Port))
}

func (a *Application) Stop() {
}

func (a *Application) redirectToIndex(c *gin.Context) {
	c.Redirect(http.StatusPermanentRedirect, "/index")
}

func (a *Application) index(c *gin.Context) {
	parameters := gin.H{}
	parameters["title"] = "Web Analyser"

	if url := c.PostForm("url"); url != "" {
		parameters["url"] = url

		log.Println("Provided URL:", url)

		doc, err := a.Fetcher.Fetch(url)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		title, ok := extractPageTitle(doc)
		if !ok {
			log.Println("Title not found")
		} else {
			log.Println("Title:", title)
		}

	}

	c.HTML(
		http.StatusOK,
		"index.html",
		parameters,
	)
}

func (a *Application) ping(c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  http.StatusOK,
			"message": "Endpoint is working",
		},
	)
}
