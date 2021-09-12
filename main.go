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

// TODO - Make UI look somewhat bearable

// TODO - update readme & github repo
// TODO - create google document

// TODO - deploy to Heroku

package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"web-analyser/fetcher"

	"github.com/gin-gonic/gin"
)

var app application

func init() {
	// Default port.
	port := 8080

	// Get port from environment.
	portString := os.Getenv("APP_PORT")
	if portString != "" {
		appPort, err := strconv.Atoi(portString)
		if err != nil {
			log.Printf("[WARNING] Invalid $APP_PORT environment variable defined ':%v', "+
				"switching to default port ':8080': %v\n", portString, err)
		} else {
			port = appPort
		}
	}

	// Get app mode from environment.
	mode := os.Getenv("APP_MODE")
	if mode == "" {
		mode = gin.DebugMode
	}

	// Initialize application.
	app.Init(port, mode, fetcher.NewURLFetcher())
}

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		err := app.Start()
		if err != nil {
			log.Fatalln("Failed to start application:", err)
		}
	}()

	s := <-signals

	log.Println("[INFO] Attempting to shutdown application due to os signal:", s.String())

	err := app.Stop()
	if err != nil {
		log.Fatalln("Failed to shutdown application gracefully", err)
	}
}
