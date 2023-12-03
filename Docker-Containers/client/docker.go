package client

import (
	"Docker-Containers/dto"
	"bufio"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

var scalableServices = []string{"hotel", "search", "user-reservation"}

func GetStats() (dto.ContainersStatsDto, error) {

	var containersStats dto.ContainersStatsDto

	command := exec.Command("docker", "stats", "--no-stream", "--format", "{{json .}}")

	stdout, err := command.StdoutPipe()
	if err != nil {
		return dto.ContainersStatsDto{}, err
	}

	err = command.Start()
	if err != nil {
		return dto.ContainersStatsDto{}, err
	}

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		var containerStats dto.ContainerStatsDto

		err = json.Unmarshal(scanner.Bytes(), &containerStats)
		if err != nil {
			return dto.ContainersStatsDto{}, err
		}

		containersStats = append(containersStats, containerStats)
	}

	err = scanner.Err()
	if err != nil {
		return dto.ContainersStatsDto{}, err
	}

	err = command.Wait()
	if err != nil {
		return dto.ContainersStatsDto{}, err
	}

	return containersStats, nil
}

func GetStatsByService(service string) (dto.ContainersStatsDto, error) {

	if !serviceExists(service) {
		return dto.ContainersStatsDto{}, errors.New("service does not exist")
	}

	var containersStats dto.ContainersStatsDto

	containers, err := getContainersIdByService(service)
	if err != nil {
		return dto.ContainersStatsDto{}, err
	}

	cmdArgs := append([]string{"stats", "--no-stream", "--format", "{{json .}}"}, containers...)

	command := exec.Command("docker", cmdArgs...)

	stdout, err := command.StdoutPipe()
	if err != nil {
		return dto.ContainersStatsDto{}, err
	}

	err = command.Start()
	if err != nil {
		return dto.ContainersStatsDto{}, err
	}

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		var containerStats dto.ContainerStatsDto

		err = json.Unmarshal(scanner.Bytes(), &containerStats)
		if err != nil {
			return dto.ContainersStatsDto{}, err
		}

		containersStats = append(containersStats, containerStats)
	}

	err = scanner.Err()
	if err != nil {
		return dto.ContainersStatsDto{}, err
	}

	err = command.Wait()
	if err != nil {
		return dto.ContainersStatsDto{}, err
	}

	return containersStats, nil
}

func getContainersIdByService(service string) ([]string, error) {

	command := exec.Command("docker-compose", "-f", "../docker-compose.yml", "ps", "-q", service)
	output, err := command.Output()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	ids := strings.TrimSpace(string(output))

	idsArray := strings.Split(ids, "\n")

	return idsArray, nil
}

func serviceExists(service string) bool {

	command := exec.Command("docker-compose", "-f", "../docker-compose.yml", "config", "--services")
	output, err := command.Output()
	if err != nil {
		log.Error(err)
		return false
	}

	services := strings.TrimSpace(string(output))
	servicesArray := strings.Split(services, "\n")

	for _, serv := range servicesArray {

		if serv == service {
			return true
		}
	}

	return false
}

func serviceScalable(service string) bool {

	for _, serv := range scalableServices {
		if serv == service {
			return true
		}
	}

	return false
}
