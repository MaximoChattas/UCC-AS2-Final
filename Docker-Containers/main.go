package main

import (
	"Docker-Containers/app"
	"Docker-Containers/docker"
)

func main() {
	docker.StartClient()
	app.StartRoute()
}
