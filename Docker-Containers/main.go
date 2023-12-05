package main

import (
	"Docker-Containers/app"
	"Docker-Containers/client"
)

func main() {

	services := client.GetScalableServices()

	for _, service := range services {

		go client.AutoScale(service)

	}

	app.StartRoute()
}
