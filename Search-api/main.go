package main

import (
	"Search/app"
	"Search/queue"
	"Search/solr"
	"sync"
)

func main() {

	queue.InitQueue()
	solr.InitSolr()

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
