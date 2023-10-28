package app

import (
	"Search/controller"
	log "github.com/sirupsen/logrus"
)

func mapUrls() {

	// Add all methods and its mappings
	router.GET("/hotel/:id", controller.GetHotelById)
	router.GET("/hotel", controller.GetHotels)

	log.Info("Finishing mappings configurations")
}
