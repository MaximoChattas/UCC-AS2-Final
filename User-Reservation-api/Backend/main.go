package main

import (
	"User-Reservation/app"
	"User-Reservation/db"
	"User-Reservation/utils"
)

func main() {

	db.StartDbEngine()
	utils.InitCache()
	app.StartRoute()
}
