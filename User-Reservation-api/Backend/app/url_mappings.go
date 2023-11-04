package app

import (
	"User-Reservation/controller"
	log "github.com/sirupsen/logrus"
)

func mapUrls() {

	// Add all methods and its mappings
	router.POST("/user", controller.InsertUser)
	router.GET("/user/:id", controller.GetUserById)
	router.GET("/user", controller.GetUsers)
	router.POST("/login", controller.UserLogin)

	router.POST("/reserve", controller.InsertReservation)
	router.GET("/reservation/:id", controller.GetReservationById)
	router.GET("/reservation", controller.GetReservations)
	router.DELETE("/reservation/:id", controller.DeleteReservation)
	router.GET("/available", controller.CheckAvailability)

	log.Info("Finishing mappings configurations")
}
