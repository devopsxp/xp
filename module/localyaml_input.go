package module

import (
	"reflect"
	"sync"
	"time"

	. "github.com/devopsxp/xp/plugin"
	"github.com/devopsxp/xp/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	AddInput("localyaml", reflect.TypeOf(LocalYamlInput{}))
}

// 解析复杂ip为单个ip
func getips(data []string) ([]string, error) {
	ips := []string{}
	for _, i := range data {
		data, err := utils.ParseIps(i)
		if err != nil {
			return ips, err
		} else {
			ips = append(ips, data...)
		}
	}
	return ips, nil
}

type LocalYaml struct {
	data map[string]interface{}
}

func (l *LocalYaml) Get() {
	l.data = viper.AllSettings()
}

type LocalYamlInput struct {
	LifeCycle
	status     StatusPlugin
	yaml       LocalYaml
	connecheck map[string]string
	lock       sync.RWMutex
}

func (l *LocalYamlInput) Receive() *Message {
	l.yaml.Get()

	if l.status != Started {
		log.Warnln("LocalYaml input plugin is not running,input nothing.")
		return nil
	}

	return Builder().WithInit().WithCheck(l.connecheck).WithItemInterface(l.yaml.data).Build()
}

func (l *LocalYamlInput) SetConnectStatus(ip, status string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.connecheck[ip] = status
}

func (l *LocalYamlInput) Start() {
	l.status = Started
	log.Debugln("LocalYamlInput plugin started.")

	// Check all ips
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
				l.SetConnectStatus(ip, "success")
			} else {
				log.Infof("%d: Ssh check %s failed 耗时：%v", num, ip, time.Now().Sub(now))
				l.SetConnectStatus(ip, "failed")
			}
			<-checkchan
		}(i, n)
	}

	wg.Wait()
}

// LocalYamlInput的Init函数实现
func (l *LocalYamlInput) Init() {
	l.yaml.data = make(map[string]interface{})
	l.connecheck = make(map[string]string)
	l.name = "LocalYaml Input"
}
