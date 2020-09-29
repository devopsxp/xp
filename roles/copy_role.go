package roles

import (
	"errors"
	"fmt"
	"strings"

	. "github.com/devopsxp/xp/plugin"
	"github.com/devopsxp/xp/utils"
	log "github.com/sirupsen/logrus"
)

type CopyRole struct {
	RoleLC
	src   string // 源地址
	dest  string // 目的地址
	msg   *Message
	logs  map[string]string
	items []string
}

// 准备数据
// @Param stage 阶段标记
// @Param user 远端执行用户
// @Param host 目标主机
// @Param vars 动态参数
// @Param configs 执行模块内容
// @Param msg 消息结构体
func (r *CopyRole) Init(stage, user, host string, vars map[string]interface{}, data map[interface{}]interface{}, msg *Message) error {
	// 获取name
	if name, ok := data["name"]; !ok {
		return errors.New("config 无 name字段")
	} else {
		r.name = name.(string)
	}

	// 获取stage
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

	copyData := data["copy"].(map[interface{}]interface{})
	// 获取原始shell命令
	r.src = copyData["src"].(string)
	r.dest = copyData["dest"].(string)

	if item, ok := data["with_items"]; ok {
		for _, it := range item.([]interface{}) {
			r.items = append(r.items, fmt.Sprintf("%v", it))
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

func (r *CopyRole) Run() error {
	var err error
	if r.items == nil {
		err := utils.New(r.host, r.remote_user, "", 22).SftpUploadToRemote(r.src, r.dest)
		if err != nil {
			log.WithFields(log.Fields{
				"Host":  r.host,
				"Name":  r.name,
				"src":   r.src,
				"dest":  r.dest,
				"Stage": r.stage,
				"User":  r.remote_user,
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
				"src":   r.src,
				"dest":  r.dest,
				"Stage": r.stage,
				"User":  r.remote_user,
			}).Infof("success upload file %s", r.dest)
			r.logs[fmt.Sprintf("%s %s %s", r.stage, r.host, r.name)] = fmt.Sprintf("success upload file %s", r.dest)
		}
	} else {
		for _, it := range r.items {
			// 补充go template基本语法
			// 注意：只针对with_items数组类型
			src, err := utils.ApplyTemplate(r.src, map[string]interface{}{"item": it})
			if err != nil {
				log.Errorf("src %s error: %v", src, err)
				panic(err)
			}
			dest, err := utils.ApplyTemplate(r.dest, map[string]interface{}{"item": it})
			err = utils.New(r.host, r.remote_user, "", 22).SftpUploadToRemote(src, dest)
			if err != nil {
				log.WithFields(log.Fields{
					"Host":  r.host,
					"Name":  r.name,
					"src":   src,
					"dest":  dest,
					"Stage": r.stage,
					"User":  r.remote_user,
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
					"src":   src,
					"dest":  dest,
					"Stage": r.stage,
					"User":  r.remote_user,
				}).Infof("success upload file %s", dest)
				r.logs[fmt.Sprintf("%s %s %s", r.stage, r.host, r.name)] = fmt.Sprintf("success upload file %s", dest)
			}
		}
	}
OVER:
	return err
}

// 处理返回日志
func (r *CopyRole) After() {
	r.msg.CallBack[fmt.Sprintf("%s-%s-%s", r.host, r.stage, r.name)] = r.logs
}
