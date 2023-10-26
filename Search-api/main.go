package main

import (
	"Search/queue"
	"sync"
)

func main() {
	queue.InitQueue()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		queue.Consume()
	}()

	// Wait for the goroutine to finish
	wg.Wait()
}
