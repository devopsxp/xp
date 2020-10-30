package module

import (
	"os"
	"reflect"
	"runtime"

	. "github.com/devopsxp/xp/plugin"
	"github.com/devopsxp/xp/roles"
	log "github.com/sirupsen/logrus"
)

func init() {
	// 初始化shell filter插件映射关系表
	AddFilter("shell", reflect.TypeOf(ShellFilter{}))
}

// shell 命令运行filter插件
type ShellFilter struct {
	LifeCycle
	status    StatusPlugin
	terminial bool
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

	// 解析terminial
	if terminial, ok := msgs.Data.Items["terminial"]; ok {
		s.terminial = terminial.(bool)
	}

	log.Debugf("Config %v\n", configs)
	var remote_user string
	if user, ok := msgs.Data.Items["remote_user"]; ok {
		remote_user = user.(string)
	} else {
		// 默认root用户
		remote_user = "root"
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
				// 判断stage是否允许执行
				if roles.IsRolesAllow(stage.(string), rolesData) {
					// 3. TODO: 解析yaml中shell的模块，然后进行匹配
					err := roles.NewShellRole(roles.NewRoleArgs(stage.(string), remote_user, host, vars, configs, msgs, nil, s.terminial))
					if err != nil {
						log.Errorln(err.Error())
					}
				}
			}
		}
	}

	return msgs
}

func (s *ShellFilter) Init(data interface{}) {
	s.name = "Shell Filter"
	s.status = Started
	if data != nil {
		if isterminial, ok := data.(map[string]interface{})["terminial"]; ok {
			s.terminial = isterminial.(bool)
		}
	}
}
