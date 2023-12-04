package app

import (
	"Docker-Containers/controller"
	log "github.com/sirupsen/logrus"
)

func mapUrls() {

	// Add all methods and its mappings
	router.GET("/stats", controller.GetStats)
	router.GET("/stats/:service", controller.GetStatsByService)
	router.POST("/scale/:service", controller.ScaleService)

	log.Info("Finishing mappings configurations")
}
