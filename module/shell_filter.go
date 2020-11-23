package module

import (
	"os"
	"reflect"
	"runtime"

	. "github.com/devopsxp/xp/plugin"
	"github.com/devopsxp/xp/roles"
	"github.com/devopsxp/xp/utils"
	log "github.com/sirupsen/logrus"
)

func init() {
	// 初始化shell filter插件映射关系表
	AddFilter("shell", reflect.TypeOf(ShellFilter{}))
}

// shell 命令运行filter插件
type ShellFilter struct {
	LifeCycle
	status StatusPlugin
}

func (s *ShellFilter) Process(msgs *Message) *Message {
	if s.status != Started {
		log.Warnln("Shell filter plugin is not running,filter nothing.")
		return msgs
	}

	// TODO:
	// 1. 封装config shell|copy|template等操作
	log.Infof("ShellFilter Filter 插件开始执行目标主机Config Playbook，并发数： %d", runtime.NumCPU())

	// 解析yaml结果
	log.Debugf("解析yaml结果 Check %v\n", msgs.Data.Check)
	// 1. 解析stage步骤
	var stages []interface{}
	if sg, ok := msgs.Data.Items["stage"]; ok {
		stages = sg.([]interface{})
	}
	log.Debugf("Stage %v\n", stages)

	var configs []interface{}
	if cf, ok := msgs.Data.Items["config"]; ok {
		configs = cf.([]interface{})
	} else {
		log.Errorln("未配置config模块，退出！")
		os.Exit(1)
	}

	// 解析include目录文件加载Yaml
	// include: /tmp/d.yaml
	// name: "name"
	for _, cc := range configs {
		rt, ok := roles.ParseRoleType(cc.(map[interface{}]interface{}))
		if ok && rt == roles.IncludeType {
			// 获取include路径
			includePath, ok := cc.(map[interface{}]interface{})["include"]
			if ok {
				log.Infof("匹配到 include 配置[%s] %s", cc.(map[interface{}]interface{})["name"], includePath)
				// include 配置格式为： [map[interface{}]interface{}]
				iData, err := utils.ReadYamlConfig(includePath.(string))
				if err != nil {
					log.Errorf("读取include yaml文件错误：%s", err.Error())
					os.Exit(1)
				}

				switch iData.(type) { //v表示b1 接口转换成Bag对象的值
				case []interface{}:
					configs = append(configs, iData.([]interface{})...)
				case map[interface{}]interface{}:
					configs = append(configs, iData.(map[interface{}]interface{}))
				default:
					log.Warnf("Include Yaml文件格式不能匹配 %v", iData)
				}
			}
		}
	}

	log.Debugln("configs", configs)

	log.Debugf("Config %v\n", configs)
	var (
		remote_user, remote_pwd string
		remote_port             int
	)

	if user, ok := msgs.Data.Items["remote_user"]; ok {
		remote_user = user.(string)
	} else {
		// 默认root用户
		remote_user = "root"
	}

	if pwd, ok := msgs.Data.Items["remote_pwd"]; ok {
		remote_pwd = pwd.(string)
	} else {
		// 默认root用户
		remote_pwd = ""
	}

	if port, ok := msgs.Data.Items["remote_port"]; ok {
		remote_port = port.(int)
	} else {
		// 默认root用户
		remote_port = 22
	}

	// 全局动态变量
	var vars map[string]interface{}
	if vv, ok := msgs.Data.Items["vars"]; ok {
		vars = vv.(map[string]interface{})
	} else {
		vars = make(map[string]interface{})
	}

	rolesData := msgs.Data.Items["roles"].([]interface{})
	// 2. 根据stage进行解析
	for host, status := range msgs.Data.Check {
		if status == "failed" {
			log.Debugf("host %s is failed, next.\n", host)
		} else {
			for _, stage := range stages {
				if roles.IsRolesAllow(stage.(string), rolesData) {
					// 3. TODO: 解析yaml中shell的模块，然后进行匹配
					err := roles.NewShellRole(roles.NewRoleArgs(stage.(string), remote_user, remote_pwd, host, vars, configs, msgs, nil, remote_port))
					if err != nil {
						log.Debugln(err.Error())
						os.Exit(1)
					}
				}
			}

			// execChan := make(chan string, runtime.NumCPU())
			// var w sync.WaitGroup
			// for _, stage := range stages {
			// w.Add(1)
			// execChan <- stage.(string)
			// go func() {
			// 	defer w.Done()
			// 	// 判断stage是否允许执行
			// 	if roles.IsRolesAllow(stage.(string), rolesData) {
			// 		// 3. TODO: 解析yaml中shell的模块，然后进行匹配
			// 		err := roles.NewShellRole(roles.NewRoleArgs(stage.(string), remote_user, host, vars, configs, msgs, nil))
			// 		if err != nil {
			// 			log.Debugln(err.Error())
			// 			os.Exit(1)
			// 		}
			// 	}
			// 	<-execChan
			// }()
			// }
			// w.Wait()
		}
	}

	return msgs
}

func (s *ShellFilter) Init(data interface{}) {
	s.name = "Shell Filter"
	s.status = Started
}
