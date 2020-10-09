package roles

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	. "github.com/devopsxp/xp/plugin"
	"github.com/devopsxp/xp/utils"
	log "github.com/sirupsen/logrus"
)

type TemplateRole struct {
	RoleLC
	src  string // 源地址
	dest string // 目的地址
	msg  *Message
	logs map[string]string
}

// 准备数据
// @Param stage 阶段标记
// @Param user 远端执行用户
// @Param host 目标主机
// @Param vars 动态参数
// @Param configs 执行模块内容
// @Param msg 消息结构体
func (r *TemplateRole) Init(stage, user, host string, vars map[string]interface{}, data map[interface{}]interface{}, msg *Message) error {
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

	copyData := data["template"].(map[interface{}]interface{})
	// 获取原始shell命令
	r.src = copyData["src"].(string)
	r.dest = copyData["dest"].(string)

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

// 操作流程：1. 获取vars和template 2. 解析 3. 上传
func (r *TemplateRole) Run() error {
	// 读取j2文件
	templateFile, err := ioutil.ReadFile(r.src)
	if err != nil {
		return err
	}

	destFile, err := utils.ApplyTemplate(string(templateFile), r.vars)
	if err != nil {
		return err
	}

	log.Debugf("template is %s", destFile)

	err = utils.New(r.host, r.remote_user, "", 22).SftpUploadTemplateString(destFile, r.dest)
	if err != nil {
		log.WithFields(log.Fields{
			"Host":     r.host,
			"Name":     r.name,
			"template": r.src,
			"dest":     r.dest,
			"Stage":    r.stage,
			"User":     r.remote_user,
			"耗时":       time.Now().Sub(r.starttime),
		}).Errorln(err.Error())
		r.logs[fmt.Sprintf("%s %s %s", r.stage, r.host, r.name)] = err.Error()
		if strings.Contains(err.Error(), "ssh:") {
			err = errors.New("ssh: handshake failed")
			return err
		}
	} else {
		log.WithFields(log.Fields{
			"Host":     r.host,
			"Name":     r.name,
			"template": r.src,
			"dest":     r.dest,
			"Stage":    r.stage,
			"User":     r.remote_user,
			"耗时":       time.Now().Sub(r.starttime),
		}).Infof("模板上传成功 %s", r.dest)
		r.logs[fmt.Sprintf("%s %s %s", r.stage, r.host, r.name)] = fmt.Sprintf("模板上传成功 %s", r.dest)
	}

	return nil
}

// 处理返回日志
func (r *TemplateRole) After() {
	stoptime := time.Now()
	r.logs["耗时"] = fmt.Sprintf("%v", stoptime.Sub(r.starttime))
	r.msg.CallBack[fmt.Sprintf("%s-%s-%s", r.host, r.stage, r.name)] = r.logs
}
