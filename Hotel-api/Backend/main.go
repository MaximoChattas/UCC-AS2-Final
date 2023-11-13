package main

import (
	"Hotel/app"
	"Hotel/db"
	"Hotel/queue"
)

func main() {

	db.Init_db()
	queue.QueueProducer.InitQueue()
	app.StartRoute()
}
