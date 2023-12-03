package app

import (
	"Docker-Containers/controller"
	log "github.com/sirupsen/logrus"
)

func mapUrls() {

	// Add all methods and its mappings
	router.GET("/stats", controller.GetStats)
	router.GET("/stats/:service", controller.GetStatsByService)

	log.Info("Finishing mappings configurations")
}
