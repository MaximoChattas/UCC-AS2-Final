package docker

import (
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

var DockerClient *client.Client

func StartClient() {

	var err error
	DockerClient, err = client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		log.Error("Failed to start Docker Client")
		log.Fatal(err)
		return
	}

	log.Info("Docker client started successfully")

}
