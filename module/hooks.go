package module

import (
	. "github.com/devopsxp/xp/plugin"
	log "github.com/sirupsen/logrus"
)

// 对外接口
type Alert interface {
	Send(*Message)
}

// 对内实现
type HookMethod interface {
	SpecificSend() (string, error)
	IsCurrent() bool // 判断当前是否可以发送告警信
}

// hook适配器
func NewHookAdapter(in HookMethod) *hook {
	return &hook{
		HookMethod: in,
	}
}

// 转换对外接口调用对内接口
// output 钩子结构体
// 负责处理发送
type hook struct {
	HookMethod
	Type   string
	Target string
}

func (h *hook) Send(msg *Message) {
	switch h.Type {
	case "console":
		log.Printf("console hook send %v\n", msg.Data.Check)
		for k, v := range msg.CallBack {
			log.Warnln(k, v)
		}
	default:
		log.Debugln("email hook send")
		status := h.IsCurrent()
		if !status {
			log.Warnln("不在发送时间，停止发送")
		} else {
			rs, err := h.SpecificSend()
			if err != nil {
				log.Errorln(err)
			} else {
				log.Warnln(rs)
			}
		}
	}
}

func (h *hook) SetType(t string) *hook {
	h.Type = t
	return h
}

func (h *hook) SetTarget(target string) *hook {
	h.Target = target
	return h
}
