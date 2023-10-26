package main

import (
	"Hotel/app"
	"Hotel/db"
	"Hotel/queue"
)

func main() {

	db.Init_db()
	queue.InitQueue()
	app.StartRoute()
}
