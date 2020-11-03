package roles

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	. "github.com/devopsxp/xp/plugin"
	log "github.com/sirupsen/logrus"
)

var rf *RoleFactory

func init() {
	rf = &RoleFactory{}
}

// role插件工厂对象，实现Role接口
type RoleFactory struct{}

// 读取配置，通过反射机制进行对象实例化
func (r *RoleFactory) Create(conf RoleType) (RolePlugin, error) {
	t, ok := roleNames[conf]
	if !ok {
		return nil, errors.New(fmt.Sprintf("not such role plugin: %s", conf))
	}

	// 根据reflect创建对象
	p := reflect.New(t).Interface().(RolePlugin)
	return p, nil
}

// 判断stage是否在roles执行范围
func IsRolesAllow(stage string, roles []interface{}) bool {
	for _, role := range roles {
		if stage == role.(string) {
			log.Debugf("%s stage 允许执行", stage)
			return true
		}
	}
	log.Debugf("%s stage 不允许执行", stage)
	return false
}

// 自动匹配roleNames对象和目标对象的匹配度
// 实现yaml即不用新增type字段又能自动匹配Role插件
func ParseRoleType(config map[interface{}]interface{}) (rt RoleType, isok bool) {
	isok = false
	// 遍历字段，匹配模块
	for k, _ := range config {
		// 与roleNames资源池匹配
		if _, ok := roleNames[RoleType(k.(string))]; ok {
			log.Debugf("匹配到 Role 资源池对象 roleNames %s", k)
			rt = RoleType(k.(string))
			isok = true
		}
	}
	return
}

// @Params stage 当前stage
// @Params user 远端目标机执行账户
// @Params host 执行目标机
// @Params vars 公共环境变量
// @Params configs 执行playbook
// @Params currentConfig 当前config
// @Params msg pipeline消息传递 TODO: context替换，传递上下文
// @Params hook 自定义config执行完的钩子函数
// @Params isTerminial 是否terminial执行shell
type RoleArgs struct {
	stage, user, host string
	vars              map[string]interface{}
	configs           []interface{}
	currentConfig     map[interface{}]interface{}
	msg               *Message
	hook              *Hook
	isTerminial       bool
}

func NewRoleArgs(stage, user, host string, vars map[string]interface{}, configs []interface{}, msg *Message, hook *Hook, terminial bool) *RoleArgs {
	return &RoleArgs{
		stage:       stage,
		user:        user,
		host:        host,
		vars:        vars,
		configs:     configs,
		msg:         msg,
		hook:        hook,
		isTerminial: terminial,
	}
}

// 处理config module适配
func NewShellRole(args *RoleArgs) error {
	// 判断hook是否配置
	if args.hook == nil {
		// 准备常规 hook
		args.hook = &Hook{
			isHook:   true,
			hookArgs: []string{"test", "hook", "example"},
			hookFunc: func(args ...[]string) error {
				log.Debugf("钩子函数测试Demo: %v", args)
				return nil
			},
		}
	}

	for n, config := range args.configs {
		// 设置当前config
		// 静止并行执行
		args.currentConfig = config.(map[interface{}]interface{})
		log.Debugf("当前步骤: %d 当前Stage: %s Config信息: %v", n, args.stage, config)
		rt, ok := ParseRoleType(config.(map[interface{}]interface{}))
		if !ok {
			return errors.New(fmt.Sprintf("未匹配到目标Role %v", config))
		}

		// 排除Include Tags
		if rt != IncludeType {
			// 根据RoleType创建对应Role类型
			role, err := rf.Create(rt)
			if err != nil {
				return err
			}

			// 初始化role
			err = role.Init(args)
			if err != nil {
				// 判断是stage不匹配还是其它错误
				if strings.Contains(err.Error(), "not equal") || strings.Contains(err.Error(), "不在可执行主机范围内") {
					log.Debugf("%s %v", args.host, err)
				} else {
					return err
				}
			} else {
				// 执行Role
				role.Pre()
				role.Before()
				// 处理重试逻辑
				if retry, ok := config.(map[interface{}]interface{})["retry"]; ok {
					for i := 0; i < retry.(int); i++ {
						err := role.Run()
						if err != nil {
							log.Warningf("重试第 %d 次，主机: %s Stage: %s User: %s 错误信息： %s", i, args.host, args.stage, args.user, err.Error())
							if i+1 == retry {
								log.Errorf("重试次数 %d 完毕，未能执行完成，错误信息: %s", i, err.Error())
								return err
							}
							// 重试等待时间
							if retryWait, ok := config.(map[interface{}]interface{})["retryWait"]; ok {
								log.Warnf("重试等待时间: %d 秒", retryWait.(int))
								time.Sleep(time.Duration(retryWait.(int)) * time.Second)
							} else {
								log.Warnln("重试等待时间: 3 秒")
								time.Sleep(3 * time.Second)
							}
						} else {
							break
						}
					}
				} else { // 如果没有设置retry字段
					err := role.Run()
					if err != nil {
						return err
					}
				}

				role.After()

				// Role钩子函数 自定义hook
				// @Param 实现里RolePlugin接口的实例
				ishook := role.IsHook()
				if ishook {
					err = role.Hooks()
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}
