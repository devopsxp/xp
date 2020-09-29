package roles

import (
	"strings"

	. "github.com/devopsxp/xp/plugin"
	log "github.com/sirupsen/logrus"
)

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

// 执行Role
// @Param 实现里RoleLifeCycle接口的实例
// @Param hook 钩子函数（types，target） |types include email,wechat,phone,sms and so on.
func ExecRole(role RoleLifeCycle, hook func(string, string) error) error {
	role.Pre()
	role.Before()
	err := role.Run()
	if err != nil {
		return err
	}
	role.After()

	types, target, ishook := role.IsHook()
	if ishook {
		err = role.Hooks(types, target, hook)
		if err != nil {
			return err
		}
	}
	return nil
}

// 处理config module适配
func NewShellRole(stage, user, host string, vars map[string]interface{}, configs []interface{}, msg *Message) error {
	for _, config := range configs {
		// 判断是否含有shell模块
		if _, ok := config.(map[interface{}]interface{})["shell"]; ok {
			sr := &ShellRole{}
			err := sr.Init(stage, user, host, vars, config.(map[interface{}]interface{}), msg)
			if err != nil {
				// 判断是stage不匹配还是其它错误
				if strings.Contains(err.Error(), "not equal") {
					log.Debugf("%s %v", host, err)
				} else {
					return err
				}
			} else {
				err = ExecRole(sr, nil)
				if err != nil {
					return err
				}
			}
		} else if _, ok := config.(map[interface{}]interface{})["copy"]; ok {
			copys := &CopyRole{}
			err := copys.Init(stage, user, host, vars, config.(map[interface{}]interface{}), msg)
			if err != nil {
				if strings.Contains(err.Error(), "not equal") {
					log.Debugf("%s %v", host, err)
				} else {
					return err
				}
			} else {
				err = ExecRole(copys, nil)
				if err != nil {
					return err
				}
			}
		} else if _, ok := config.(map[interface{}]interface{})["template"]; ok {
			template := &TemplateRole{}
			err := template.Init(stage, user, host, vars, config.(map[interface{}]interface{}), msg)
			if err != nil {
				if strings.Contains(err.Error(), "not equal") {
					log.Debugf("%s %v", host, err)
				} else {
					return err
				}
			} else {
				err = ExecRole(template, nil)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
