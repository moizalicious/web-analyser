package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Application struct {
	router *gin.Engine
	Port   int
}

func (a *Application) Init(port int) {
	a.Port = port
	a.router = gin.Default()

	a.router.LoadHTMLGlob("templates/*.html")

	a.router.GET("/", a.redirectToIndex)
	a.router.GET("/index", a.index)
	a.router.GET("/ping", a.ping)
}

func (a *Application) Start() error {
	return a.router.Run(":" + strconv.Itoa(a.Port))
}

func (a *Application) Stop() {
}

func (a *Application) redirectToIndex(c *gin.Context) {
	c.Redirect(http.StatusPermanentRedirect, "/index")
}

func (a *Application) index(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"index.html",
		gin.H{
			"title":   "Web Analyser",
			"heading": "Hello Gin New With Params",
		},
	)
}

func (a *Application) ping(c *gin.Context) {
	type msg struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}

	c.JSON(
		http.StatusOK,
		msg{
			Status:  http.StatusOK,
			Message: "200 OK",
		},
	)
}
