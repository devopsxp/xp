package module

import (
	"fmt"
	"log/slog"
	"reflect"

	"github.com/devopsxp/xp/pkg/k8s"
	"github.com/devopsxp/xp/plugin"
)

func init() {
	// 初始化output插件映射关系表
	plugin.AddOutput("console", reflect.TypeOf(ConsoleOutput{}))
}

// Console output插件，将消息输出到控制台上
type ConsoleOutput struct {
	LifeCycle
	status plugin.StatusPlugin
}

func (c *ConsoleOutput) Send(msgs *plugin.Message) {
	if c.status != plugin.Started {
		slog.Warn("Console output is not running, output nothing.")
		return
	}

	// slog.Printf("Output:\n\tHeader: %+v, Body: %+v\n", msgs.Data.Raw, msgs.Data.Target)
	// c.SetType("console").SetTarget("stdout").Send(msgs)
	slog.Info("ConsoleOutput Output 插件开始执行目标主机，并发数： 1")

	// 全局动态变量
	var vars map[string]interface{}
	if vv, ok := msgs.Data.Items["vars"]; ok {
		vars = vv.(map[string]interface{})
	} else {
		vars = make(map[string]interface{})
	}
	// 获取hook配置 默认为console
	// 添加hooks不存在的默认配置console
	if sendtypes, ok := msgs.Data.Items["hooks"]; ok {
		if len(sendtypes.([]interface{})) > 0 {
			for _, types := range sendtypes.([]interface{}) {
				if t, ok := types.(map[interface{}]interface{})["type"]; ok {
					switch t.(string) {
					case "k8shook":
						NewHookAdapter(nil).SetType("count").Send(msgs)
						if ns, ok := msgs.Tmp["namespace"]; ok {
							if name, ok := msgs.Tmp["name"]; ok {
								slog.Info(fmt.Sprintf("Pipeline Success,清理 Namespace: %s Pod: %s", ns, name))
								err := k8s.DeletePod(ns, name)
								if err != nil {
									slog.Error(fmt.Sprintf("Pipeline 清理失败， Namespace： %s Pod: %s %s", ns, name, err.Error()))
								}
							}
						}
					case "count":
						NewHookAdapter(nil).SetType("count").Send(msgs)
					case "console":
						NewHookAdapter(nil).SetType("console").Send(msgs)
					case "email":
						email, err := NewEmail(types.(map[interface{}]interface{}), msgs, vars)
						if err != nil {
							slog.Error(err.Error())
						} else {
							NewHookAdapter(email).SetType("email").Send(msgs)
						}
					case "wechat":
						wechat, err := NewWechat(types.(map[interface{}]interface{}), msgs, vars)
						if err != nil {
							slog.Error(err.Error())
						} else {
							NewHookAdapter(wechat).SetType("wechat").Send(msgs)
						}
					default:
						slog.Warn(fmt.Sprintf("未适配该类型的hooks: %s", t.(string)))
					}
				} else {
					slog.Error("hooks 配置内容不包含[type]字段,请检查！")
					break
				}
			}
		}
	} else {
		NewHookAdapter(nil).SetType("console").Send(msgs)
	}

}

func (c *ConsoleOutput) Init(data interface{}) {
	c.name = "Console output"
	c.status = plugin.Started
}
