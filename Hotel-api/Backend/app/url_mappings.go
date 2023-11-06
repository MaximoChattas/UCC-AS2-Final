package app

import (
	"Hotel/controller"
	log "github.com/sirupsen/logrus"
)

func mapUrls() {

	// Add all methods and its mappings
	router.POST("/hotel", controller.InsertHotel)
	router.GET("/hotel/:id", controller.GetHotelById)
	router.GET("/hotel", controller.GetHotels)
	router.POST("/hotel/:id/images", controller.InsertImages)
	router.DELETE("/hotel/:id", controller.DeleteHotel)
	router.PUT("/hotel/:id", controller.UpdateHotel)
	router.GET("/image", controller.GetImageByName)

	router.POST("/amenity", controller.InsertAmenity)
	router.GET("/amenity", controller.GetAmenities)
	router.DELETE("amenity/:id", controller.DeleteAmenityById)

	log.Info("Finishing mappings configurations")
}
