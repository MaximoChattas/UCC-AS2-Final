package main

import (
	"User-Reservation/app"
	"User-Reservation/cache"
	"User-Reservation/db"
)

func main() {

	db.StartDbEngine()
	cache.InitCache()
	app.StartRoute()
}
