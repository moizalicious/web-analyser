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

// TODO - add unit tests and benchmarks
// TODO - deploy to Heroku

// TODO - Make UI look somewhat bearable
// TODO - Comment all files
// TODO - clean all the code
// TODO - update readme
// TODO - create google document

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
	port := 8080

	portString := os.Getenv("APP_PORT")
	if portString != "" {
		appPort, err := strconv.Atoi(portString)
		if err != nil {
			log.Printf("[WARNING] Invalid $APP_PORT environment variable defined ':%v', "+
				"switching to default port ':8080': %v", portString, err)
		} else {
			port = appPort
		}
	}

	mode := os.Getenv("APP_MODE")
	if mode == "" {
		mode = gin.DebugMode
	}

	app.Init(port, mode, fetcher.NewFetcher())
}

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		err := app.Start()
		if err != nil {
			log.Fatalf("Failed to start application: %v\n", err)
		}
	}()

	s := <-signals

	log.Printf("Gracefully shutting down service due to os signal '%v'\n", s)

	app.Stop()
}
