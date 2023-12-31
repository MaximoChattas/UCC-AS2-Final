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
	router.GET("/user/reservations/:id", controller.GetReservationsByUser)
	router.GET("/available", controller.CheckAvailability)

	router.POST("/amadeus", controller.InsertAmadeusMap)
	router.GET("/amadeus/:hotel_id", controller.GetAmadeusIdByHotelId)
	router.DELETE("/amadeus/:hotel_id", controller.DeleteMapping)

	log.Info("Finishing mappings configurations")
}
