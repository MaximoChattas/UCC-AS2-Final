package main

import (
	"User-Reservation/app"
	"User-Reservation/db"
)

func main() {

	db.StartDbEngine()
	app.StartRoute()
}
