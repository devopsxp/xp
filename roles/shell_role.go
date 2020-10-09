package roles

import (
	"errors"
	"fmt"
	"strings"
	"time"

	. "github.com/devopsxp/xp/plugin"
	"github.com/devopsxp/xp/utils"
	log "github.com/sirupsen/logrus"
)

type ShellRole struct {
	RoleLC
	shell string // 原生命令
	msg   *Message
	logs  map[string]string // 命令执行日志
	items []string          // 多命令集合
}

// 准备数据
// @Param stage 阶段标记
// @Param user 远端执行用户
// @Param host 目标主机
// @Param vars 动态参数
// @Param configs 执行模块内容
// @Param msg 消息结构体
func (r *ShellRole) Init(stage, user, host string, vars map[string]interface{}, data map[interface{}]interface{}, msg *Message) error {
	if current_stage, ok := data["stage"]; !ok {
		return errors.New("config 无 stage字段")
	} else {
		if stage != current_stage.(string) {
			return errors.New(fmt.Sprintf("stage not equal %s %d != %s %d", stage, len(stage), current_stage, len(current_stage.(string))))
		}
	}

	r.logs = make(map[string]string)
	r.msg = msg
	r.remote_user = user
	r.stage = stage
	r.vars = vars

	r.host = host

	// 获取原始shell命令
	r.shell = data["shell"].(string)

	// 获取name
	r.name = data["name"].(string)

	// 获取with_items迭代
	if item, ok := data["with_items"]; ok {
		for _, it := range item.([]interface{}) {
			r.items = append(r.items, it.(string))
		}
	}

	// 是否在可执行主机范围内
	isTags := false

	// 获取tags目标执行主机
	if tags, ok := data["tags"]; ok {
		for _, tag := range tags.([]interface{}) {
			if host == tag.(string) {
				isTags = true
			}
		}
	} else {
		// 没有设置tags标签，表示不限制主机执行
		isTags = true
	}

	if !isTags {
		return errors.New(fmt.Sprintf("Stage: %s Name: %s Host: %s 不在可执行主机范围内，退出！", stage, r.name, host))
	}

	return nil
}

// 执行
func (r *ShellRole) Run() error {
	var err error
	if r.items == nil {
		rs, err := utils.New(r.host, r.remote_user, "", 22).Run(r.shell)
		if err != nil {
			log.WithFields(log.Fields{
				"Host":  r.host,
				"Name":  r.name,
				"Shell": r.shell,
				"Stage": r.stage,
				"User":  r.remote_user,
				"耗时":    time.Now().Sub(r.starttime),
			}).Errorln(err.Error())
			r.logs[fmt.Sprintf("%s %s %s", r.stage, r.host, r.name)] = err.Error()
			if strings.Contains(err.Error(), "ssh:") {
				err = errors.New("ssh: handshake failed")
				goto OVER
			}
		} else {
			log.WithFields(log.Fields{
				"Host":  r.host,
				"Name":  r.name,
				"Shell": r.shell,
				"Stage": r.stage,
				"User":  r.remote_user,
				"耗时":    time.Now().Sub(r.starttime),
			}).Info(rs)
			r.logs[fmt.Sprintf("%s %s %s", r.stage, r.host, r.name)] = rs
		}
	} else {
		for _, it := range r.items {
			// 补充go template基本语法
			// 注意：只针对with_items数组类型
			cmd, err := utils.ApplyTemplate(r.shell, map[string]interface{}{"item": it})
			if err != nil {
				log.Errorf("cmd %s error: %v", cmd, err)
				panic(err)
			}
			log.Debugf("cmd is %s", cmd)
			rs, err := utils.New(r.host, r.remote_user, "", 22).Run(cmd)
			if err != nil {
				log.WithFields(log.Fields{
					"Host":  r.host,
					"Name":  r.name,
					"Shell": cmd,
					"Stage": r.stage,
					"User":  r.remote_user,
					"耗时":    time.Now().Sub(r.starttime),
				}).Errorln(err.Error())
				r.logs[fmt.Sprintf("%s %s %s", r.stage, r.host, r.name)] = err.Error()
				if strings.Contains(err.Error(), "ssh:") {
					err = errors.New("ssh: handshake failed")
					goto OVER
				}
			} else {
				log.WithFields(log.Fields{
					"Host":  r.host,
					"Name":  r.name,
					"Shell": cmd,
					"Stage": r.stage,
					"User":  r.remote_user,
					"耗时":    time.Now().Sub(r.starttime),
				}).Info(rs)
				r.logs[fmt.Sprintf("%s %s %s", r.stage, r.host, r.name)] = rs
			}
		}
	}
OVER:
	return err
}

// 处理返回日志
func (r *ShellRole) After() {
	stoptime := time.Now()
	r.logs["耗时"] = fmt.Sprintf("%v", stoptime.Sub(r.starttime))
	r.msg.CallBack[fmt.Sprintf("%s-%s-%s", r.host, r.stage, r.name)] = r.logs
}

func testhook(a, b string) error {
	log.Printf("%s %s test hook send")
	return nil
}
