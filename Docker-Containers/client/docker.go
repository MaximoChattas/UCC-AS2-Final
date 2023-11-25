package client

import (
	"Docker-Containers/dto"
	"bufio"
	"encoding/json"
	"os/exec"
)

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
