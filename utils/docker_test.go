package utils

// import (
// 	"testing"
// )

// func TestDockerListContainers(t *testing.T) {
// 	host := "127.0.0.1"
// 	port := "9999"
// 	version := "v1.24"
// 	cli := NewDockerCLI(host, port, version)
// 	rs, err := cli.ListContainers()
// 	if err != nil {
// 		t.Errorf("ListContainers() return an error: %v", err)
// 	}

// 	for _, container := range rs.([]interface{}) {
// 		if rs, ok := container.(map[string]interface{})["Id"]; ok {
// 			t.Logf("success got container Id: %s", rs.(string))
// 			break
// 		} else {
// 			t.Error("can not found id")
// 			break
// 		}
// 	}
// }

// func TestContainerTop(t *testing.T) {
// 	host := "127.0.0.1"
// 	port := "9999"
// 	version := "v1.24"
// 	cli := NewDockerCLI(host, port, version)
// 	rs, err := cli.ListContainers()
// 	if err != nil {
// 		t.Errorf("ListContainers() return an error: %v", err)
// 	}

// 	for _, container := range rs.([]interface{}) {
// 		if rs, ok := container.(map[string]interface{})["Id"]; ok {
// 			// t.Logf("success got container Id: %s", rs.(string))
// 			rs, err := cli.ContainerTop(rs.(string))
// 			if err != nil {
// 				t.Error(err)
// 				break
// 			}
// 			t.Log(rs)
// 			break
// 		} else {
// 			t.Error("can not found id")
// 			break
// 		}
// 	}
// }

// func TestContainerLogs(t *testing.T) {
// 	host := "127.0.0.1"
// 	port := "9999"
// 	version := "v1.24"
// 	cli := NewDockerCLI(host, port, version)
// 	rs, err := cli.ListContainers()
// 	if err != nil {
// 		t.Errorf("ListContainers() return an error: %v", err)
// 	}

// 	for _, container := range rs.([]interface{}) {
// 		if rs, ok := container.(map[string]interface{})["Id"]; ok {
// 			// t.Logf("success got container Id: %s", rs.(string))
// 			rs, err := cli.ContainerLogs(rs.(string))
// 			if err != nil {
// 				t.Error(err)
// 				break
// 			}
// 			t.Log(rs)
// 			break
// 		} else {
// 			t.Error("can not found id")
// 			break
// 		}
// 	}
// }

// func TestCreateContainer(t *testing.T) {
// 	host, port := "127.0.0.1", "9999"
// 	data := map[string]interface{}{
// 		"Hostname":     "",
// 		"Domainname":   "",
// 		"User":         "",
// 		"AttachStdin":  false,
// 		"AttachStdout": true,
// 		"AttachStderr": true,
// 		"Tty":          false,
// 		"OpenStdin":    false,
// 		"StdinOnce":    false,
// 		"Env": []string{
// 			"FOO=bar",
// 			"BAZ=quux",
// 		},
// 		"Cmd": []string{
// 			"date",
// 		},
// 		"Entrypoint": "",
// 		"Image":      "ubuntu",
// 		"Labels": map[string]string{
// 			"com.example.vendor":  "Acme",
// 			"com.example.license": "GPL",
// 			"com.example.version": "1.0",
// 		},
// 		"Volumes": map[string]interface{}{
// 			"/tmp": map[string]string{},
// 		},
// 		"Healthcheck": map[string]interface{}{
// 			"Test":        []string{"CMD-SHELL", "curl localhost:3000"},
// 			"Interval":    1000000000,
// 			"Timeout":     10000000000,
// 			"Retries":     10,
// 			"StartPeriod": 60000000000,
// 		},
// 		"WorkingDir":      "/",
// 		"NetworkDisabled": false,
// 		"MacAddress":      "12:34:56:78:9a:bc",
// 		"ExposedPorts": map[string]interface{}{
// 			"22/tcp": map[string]string{},
// 		},
// 		"StopSignal": "SIGTERM",
// 		"HostConfig": map[string]interface{}{
// 			"Binds":             []string{"/tmp:/tmp"},
// 			"Tmpfs":             map[string]string{"/run": "rw,noexec,nosuid,size=65536k"},
// 			"Links":             []string{}, // "redis3:redis"
// 			"Memory":            0,          // 8MB
// 			"MemorySwap":        0,
// 			"MemoryReservation": 0,
// 			"KernelMemory":      0,
// 			// "CpuPercent":         80,
// 			"CpuShares":          512,
// 			"CpuPeriod":          100000,
// 			"CpuQuota":           50000,
// 			"CpusetCpus":         "",
// 			"CpusetMems":         "",
// 			"IOMaximumBandwidth": 0,
// 			"IOMaximumIOps":      0,
// 			// "BlkioWeight":          300,
// 			// "BlkioWeightDevice":    []map[string]string{},
// 			// "BlkioDeviceReadBps":   []map[string]string{},
// 			// "BlkioDeviceReadIOps":  []map[string]string{},
// 			// "BlkioDeviceWriteBps":  []map[string]string{},
// 			// "BlkioDeviceWriteIOps": []map[string]string{},
// 			"MemorySwappiness": 60,
// 			"OomKillDisable":   false,
// 			"OomScoreAdj":      500,
// 			"PidMode":          "",
// 			"PidsLimit":        -1,
// 			"PortBindings":     map[string]interface{}{"22/tcp": []map[string]string{map[string]string{"HostPort": "11022"}}},
// 			"PublishAllPorts":  false,
// 			"Privileged":       false,
// 			"ReadonlyRootfs":   false,
// 			"Dns":              []string{"8.8.8.8"},
// 			"DnsOptions":       []string{},
// 			"DnsSearch":        []string{},
// 			"ExtraHosts":       []string{},
// 			"VolumesFrom":      []string{}, // ["parent", "other:ro"],
// 			"CapAdd":           []string{"NET_ADMIN"},
// 			"CapDrop":          []string{"MKNOD"},
// 			// "GroupAdd":         []string{"newgroup"},
// 			"RestartPolicy": map[string]interface{}{"Name": "", "MaximumRetryCount": 0},
// 			"NetworkMode":   "bridge",
// 			"Devices":       []string{},
// 			"Sysctls":       map[string]string{"net.ipv4.ip_forward": "1"},
// 			"Ulimits":       []map[string]string{},
// 			"LogConfig":     map[string]interface{}{"Type": "json-file", "Config": map[string]string{}},
// 			"SecurityOpt":   []string{},
// 			// "StorageOpt":           map[string]string{},
// 			"CgroupParent": "",
// 			"VolumeDriver": "",
// 			"ShmSize":      67108864,
// 		},
// 		// "NetworkingConfig": map[string]interface{}{
// 		// 	"EndpointsConfig": map[string]interface{}{
// 		// 		"isolated_nw": map[string]interface{}{
// 		// 			"IPAMConfig": map[string]interface{}{
// 		// 				"IPv4Address":  "172.20.30.33",
// 		// 				"IPv6Address":  "2001:db8:abcd::3033",
// 		// 				"LinkLocalIPs": []string{"169.254.34.68", "fe80::3468"},
// 		// 			},
// 		// 			"Links":   []string{"container_1", "container_2"},
// 		// 			"Aliases": []string{"server_x", "server_y"},
// 		// 		},
// 		// 	},
// 		// },
// 	}
// 	cli := NewDockerCLI(host, port, "")
// 	rs, err := cli.CreateContainer(data)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	t.Log(string(rs))
// }

// func TestContainerStart(t *testing.T) {
// 	host, port := "127.0.0.1", "9999"
// 	data := map[string]interface{}{
// 		"Binds":           []string{"/tmp:/tmp"},
// 		"PortBindings":    map[string]interface{}{"22/tcp": []map[string]string{map[string]string{"HostPort": "11022"}}},
// 		"PublishAllPorts": false,
// 		"Privileged":      false,
// 	}

// 	cli := NewDockerCLI(host, port, "")
// 	rs, err := cli.ListContainers()
// 	if err != nil {
// 		t.Errorf("ListContainers() return an error: %v", err)
// 	}

// 	for _, container := range rs.([]interface{}) {
// 		if rs, ok := container.(map[string]interface{})["Id"]; ok {
// 			// t.Logf("success got container Id: %s", rs.(string))
// 			rs, err := cli.ContainerStart(rs.(string), data)
// 			if err != nil {
// 				t.Error(err)
// 				break
// 			}

// 			t.Log(string(rs))
// 			break
// 		} else {
// 			t.Error("can not found id")
// 			break
// 		}
// 	}
// }

// func TestContainerStop(t *testing.T) {
// 	host, port := "127.0.0.1", "9999"
// 	cli := NewDockerCLI(host, port, "")
// 	rs, err := cli.ListContainers()
// 	if err != nil {
// 		t.Errorf("ListContainers() return an error: %v", err)
// 	}

// 	for _, container := range rs.([]interface{}) {
// 		if rs, ok := container.(map[string]interface{})["Id"]; ok {
// 			// t.Logf("success got container Id: %s", rs.(string))
// 			rs, err := cli.ContainerStop(rs.(string))
// 			if err != nil {
// 				t.Error(err)
// 				break
// 			}

// 			t.Log(string(rs))
// 			break
// 		} else {
// 			t.Error("can not found id")
// 			break
// 		}
// 	}
// }

// func TestContainerRestart(t *testing.T) {
// 	host, port := "127.0.0.1", "9999"

// 	cli := NewDockerCLI(host, port, "")
// 	rs, err := cli.ListContainers()
// 	if err != nil {
// 		t.Errorf("ListContainers() return an error: %v", err)
// 	}

// 	for _, container := range rs.([]interface{}) {
// 		if rs, ok := container.(map[string]interface{})["Id"]; ok {
// 			// t.Logf("success got container Id: %s", rs.(string))
// 			rs, err := cli.ContainerRestart(rs.(string))
// 			if err != nil {
// 				t.Error(err)
// 				break
// 			}

// 			t.Log(string(rs))
// 			break
// 		} else {
// 			t.Error("can not found id")
// 			break
// 		}
// 	}
// }

// func TestContainerKill(t *testing.T) {
// 	host, port := "127.0.0.1", "9999"

// 	cli := NewDockerCLI(host, port, "")
// 	rs, err := cli.ListContainers()
// 	if err != nil {
// 		t.Errorf("ListContainers() return an error: %v", err)
// 	}

// 	for _, container := range rs.([]interface{}) {
// 		if rs, ok := container.(map[string]interface{})["Id"]; ok {
// 			// t.Logf("success got container Id: %s", rs.(string))
// 			rs, err := cli.ContainerRestart(rs.(string))
// 			if err != nil {
// 				t.Error(err)
// 				break
// 			}

// 			t.Log(string(rs))
// 			break
// 		} else {
// 			t.Error("can not found id")
// 			break
// 		}
// 	}
// }
