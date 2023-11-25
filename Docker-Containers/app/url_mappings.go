package app

import (
	"Docker-Containers/controller"
	log "github.com/sirupsen/logrus"
)

func mapUrls() {

	// Add all methods and its mappings
	router.GET("/stats", controller.GetStats)

	log.Info("Finishing mappings configurations")
}
