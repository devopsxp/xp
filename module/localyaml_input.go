package module

import (
	"fmt"
	"log/slog"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/devopsxp/xp/plugin"
	"github.com/devopsxp/xp/utils"
	"github.com/spf13/viper"
)

func init() {
	plugin.AddInput("localyaml", reflect.TypeOf(LocalYamlInput{}))
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
	slog.Debug(fmt.Sprintf("所有配置： %v", l.data))
}

type LocalYamlInput struct {
	LifeCycle
	status     plugin.StatusPlugin
	yaml       LocalYaml
	connecheck map[string]string
	lock       sync.RWMutex
	fails      int // ssh连接失败数
}

func (l *LocalYamlInput) Receive() *plugin.Message {
	l.yaml.Get()

	if l.status != plugin.Started {
		slog.Warn("LocalYaml input plugin is not running,input nothing.")
		return nil
	}

	return plugin.Builder().WithInit(l.fails).WithCheck(l.connecheck).WithItemInterface(l.yaml.data).Build()
}

func (l *LocalYamlInput) SetConnectStatus(ip, status string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.connecheck[ip] = status
}

func (l *LocalYamlInput) Start() {
	l.fails = 0
	l.status = plugin.Started
	slog.Debug("LocalYamlInput plugin started.")

	// Check all ips
	ips, err := getips(viper.GetStringSlice("host"))
	if err != nil {
		panic(err)
	}

	port := viper.GetInt("remote_port")
	if port <= 0 || port < 22 {
		port = 22
	}

	// 目标主机22端口检测并发限制
	checkchan := make(chan string, 5*runtime.NumCPU())

	var wg sync.WaitGroup

	slog.Info("******************************************************** TASK [LocalYamlCheck : 主机状态检测] ********************************************************")
	slog.Info("LocalYaml Input 插件开始执行ssh目标主机状态扫描", "并发数", 5*runtime.NumCPU())
	for n, i := range ips {
		wg.Add(1)
		checkchan <- i
		go func(ip string, num int) {
			defer wg.Done()
			now := time.Now()
			if utils.ScanPort(ip, port) {
				slog.Info(fmt.Sprintf("%d: Ssh check %s:%d success 耗时: %v", num, ip, port, time.Now().Sub(now)))
				l.SetConnectStatus(ip, "success")
			} else {
				slog.Info(fmt.Sprintf("%d: Ssh check %s:%d failed 耗时：%v", num, ip, port, time.Now().Sub(now)))
				l.fails += 1
				l.SetConnectStatus(ip, "failed")
			}
			<-checkchan
		}(i, n)

		if n%10 == 0 {
			slog.Info(fmt.Sprintf("已完成 %d 主机连接测试, 当前GoRoutine数量: %d", n, runtime.NumGoroutine()))
		}
	}

	wg.Wait()
}

// LocalYamlInput的Init函数实现
func (l *LocalYamlInput) Init(data interface{}) {
	l.yaml.data = make(map[string]interface{})
	l.connecheck = make(map[string]string)
	l.name = "LocalYaml Input"
}
