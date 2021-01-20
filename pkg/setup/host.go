// 通用设置模块 获取主机基本配置
package setup

import (
	"encoding/json"
	"net"

	"github.com/shirou/gopsutil/host"
)

type IP struct {
	Device   string `json:"device"`
	Ip       string `json:"ip"`
	Mac      string `json:"mac"`
	Bordcast string `json:"bordcast"`
	Status   string `json:"status"`
}

type Host struct {
	info host.InfoStat
	Ip   []IP `json:"ip"`
}

func (h *Host) getHost() error {
	data, err := host.Info()
	if err != nil {
		return err
	}

	h.info = data
	return nil
}

func (h *Host) getIps() error {
	if h.Ip == nil {
		h.Ip = []IP{}
	}

	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return err
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := IP{
					Ip: ipnet.IP.String(),
				}
				h.Ip = append(h.Ip, ipnet.IP.String())
			}
		}
	}
	return nil
}

func (h *Host) String() (string, error) {
	err := h.getHost()
	if err != nil {
		return "", err
	}

	err = h.getIps()
	if err != nil {
		return "", err
	}
	s, _ := json.Marshal(h)
	return string(s), nil
}
