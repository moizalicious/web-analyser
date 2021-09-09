package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Application struct {
	Router *gin.Engine
	Port   int
}

func (a *Application) Init(port int) {
	a.Port = port
	a.Router = gin.Default()

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
	log.Printf("URL Value From Form: %v", c.PostForm("url"))

	parameters := gin.H{}
	parameters["title"] = "Web Analyser"

	if url := c.PostForm("url"); url != "" {
		parameters["url"] = url
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
