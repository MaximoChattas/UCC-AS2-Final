package main

import (
	"Search/app"
	"Search/controller"
	"Search/queue"
	"Search/solr"
	"sync"
)

func main() {

	queue.InitQueue()
	solr.InitSolr()

	controller.InsertData()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		queue.Consume()
	}()

	app.StartRoute()

	// Wait for the goroutine to finish
	wg.Wait()
}
