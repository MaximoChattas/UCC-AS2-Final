package main

import (
	"Hotel/app"
	"Hotel/db"
)

func main() {

	db.Init_db()
	app.StartRoute()
}
