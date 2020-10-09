package utils

import (
	"encoding/json"
	"fmt"
)

func DockerListContainers(host, port, version string) (result interface{}, err error) {
	url := fmt.Sprintf("http://%s:%s/%s/containers/json?all=1", host, port, version)
	fmt.Println(url)
	data, err := Get(url)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return
	}

	return
}
