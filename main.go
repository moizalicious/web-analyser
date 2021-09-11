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
	"os"
	"os/signal"
	"syscall"
	"web-analyser/fetcher"
)

var App Application

func init() {
	// TODO - add release mode and debug mode
	App.Init(8080, fetcher.NewFetcher())
}

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		err := App.Start()
		if err != nil {
			log.Fatalf("Failed to start application: error - %v\n", err)
		}
	}()

	s := <-signals

	log.Printf("Gracefully shutting down service due to os signal '%v'\n", s)

	App.Stop()
}
