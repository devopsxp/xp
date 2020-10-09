package utils

import (
	"encoding/json"
	"fmt"
)

type docker struct {
	host    string
	port    string
	version string
}

func NewDockerCLI(host, port, version string) *docker {
	if version == "" {
		version = "v1.24"
	}

	return &docker{
		host:    host,
		port:    port,
		version: version,
	}
}

func (d *docker) ListContainers() (result interface{}, err error) {
	url := fmt.Sprintf("http://%s:%s/%s/containers/json?all=1", d.host, d.port, d.version)
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

// TODO: this is bad function
func (d *docker) CreateContainer() (result interface{}, err error) {
	url := fmt.Sprintf("http://%s:%s/%s/containers/create", d.host, d.port, d.version)
	fmt.Println(url)
	data, err := Post(url, "", "")
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return
	}

	return
}
