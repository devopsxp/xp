package module

import (
	"sync"
	"time"

	. "github.com/devopsxp/xp/plugin"
	"github.com/devopsxp/xp/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// 1. 获取cli参数
// 2. 想办法传进来
// 3. 拼装成Message.Data map[string]interface{}
// 4. 执行pipeline

type CliInput struct {
	LifeCycle
	status       StatusPlugin
	connectcheck map[string]string
	lock         sync.RWMutex
}

func (c *CliInput) Receive(data interface{}) *Message {
	if c.status != Started {
		log.Warnln("LocalYaml input plugin is not running,input nothing.")
		return nil
	}

	return Builder().WithInit().WithCheck(c.connectcheck).WithItemInterface(data.(map[string]interface{})).Build()
}

func (c *CliInput) SetConnectStatus(ip, status string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.connectcheck[ip] = status
}

func (c *CliInput) Start() {
	c.status = Started
	log.Debugln("LocalYamlInput plugin started.")

	// Check all ipsl.yaml.
	// TODO: error 没有viper取配置了
	ips, err := getips(viper.GetStringSlice("host"))
	if err != nil {
		panic(err)
	}

	// 目标主机22端口检测并发限制
	checkchan := make(chan string, 10)

	var wg sync.WaitGroup

	log.Info("LocalYaml Input 插件开始执行ssh目标主机状态扫描，并发数： 10")
	for n, i := range ips {
		wg.Add(1)
		go func(ip string, num int) {
			defer wg.Done()
			checkchan <- ip
			now := time.Now()
			if utils.ScanPort(ip, "22") {
				log.Infof("%d: Ssh check %s success 耗时: %v", num, ip, time.Now().Sub(now))
				c.SetConnectStatus(ip, "success")
			} else {
				log.Infof("%d: Ssh check %s failed 耗时：%v", num, ip, time.Now().Sub(now))
				c.SetConnectStatus(ip, "failed")
			}
			<-checkchan
		}(i, n)
	}

	wg.Wait()
}

// LocalYamlInput的Init函数实现
func (c *CliInput) Init() {
	c.connectcheck = make(map[string]string)
	c.name = "Cli Input"
}
