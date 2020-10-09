package utils

import (
	"testing"
)

func TestDockerListContainers(t *testing.T) {
	host := "127.0.0.1"
	port := "9999"
	version := "v1.24"
	rs, err := DockerListContainers(host, port, version)
	if err != nil {
		t.Errorf("DockerListContainers() return an error: %v", err)
	}

	for _, container := range rs.([]interface{}) {
		if rs, ok := container.(map[string]interface{})["Id"]; ok {
			t.Logf("success got container Id: %s", rs.(string))
			break
		} else {
			t.Error("can not found id")
			break
		}
	}
}
