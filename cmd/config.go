/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const yaml string = `host: # 目标主机
  - 127.0.0.1
  - 192.168.0.1-255
remote_user: root 
roles: # 执行具体stage
  - build
  - test
vars: # 环境变量
  Status: true
  title: "模板测试"
  serviceName: "xp"
  pipelineName: template test
  pipelineId: no1
  pipelineUrl: http://www.baidu.com
  branch: master
  stage: build
  tags: test
  info: test only
  alerts:
    - status: firing
      generatorURL: http://www.google.com
      startsAt: 2020-08-01
      endsAt: 2020-09-01
      annotations:
        current_value: 85
      labels:
        severity: warning
        node: 127.0.0.1
        threshold_value: 80
    - status: ok
      generatorURL: http://www.google2.com
      startsAt: 2020-18-01
      endsAt: 2020-29-01
      annotations:
        current_value: 99
      labels:
        severity: ok
        node: 127.0.0.1
        threshold_value: 88
hooks: # 消息发送 全局一个 可以一次发送所有结果到多个渠道
  - type: console # 类型，支持：console|email|wechat|phone|sms
  - type: email
    alias: xp战队 # 邮箱昵称
    email_user: "xp@xp.com"
    email_pwd: "******"
    email_smtp: "smtp.exmail.qq.com"
    email_smtp_port: "465"
    email_to:
      - xp@xp.com
    template: # 告警模板
      title: "告警title"
      text: "{{.Status}} {{.title}}
      {{range $key,$value := .logs}}
        {{range $k,$v := $value}}
        ------ {{$k}} {{$v}}<br/>
        {{end}}
      {{end}}"
      path: template.service.j2
      vars: # 内置固定参数
        Status: true
        title: "模板测试"
        serviceName: "xp"
  - type: wechat
    address:
      - https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=******
    template: # 告警模板
      text: '
# 告警模板

{{if .Status}} <font color=\"info\">[成功] </font>{{ .title}} {{- else}}<font color=\"warning\">[失败] </font> {{ .title}}  {{- end}} 

> 服务名称：<font color="comment">{{.serviceName}}</font> 
> 流水线名称：<font color="comment">{{.pipelineName}}</font> 
> 流水线ID：<font color="comment">{{.pipelineId}}</font> 
> 流水线地址：<font color="comment">{{.pipelineUrl}}</font> 
> 构建分支：<font color="comment">{{.branch}}</font> 
> 阶段：<font color="comment">{{.stage}}</font> 
> 构建人：<font color="comment">{{.user}}</font> 
> COMMIT: <font color="comment">{{.info}}</font> 
> COMMITID：<font color="comment">{{.tags}}</font>

# {{.title}}

{{range $key,$value := .logs}}
{{range $k,$v := $value}}
> {{$k}} {{$v}}
{{end}}
{{end}}'
      path: template.service.j2
      vars: # 内置固定参数
        title: "模板测试"
        serviceName: "xp"
  - type: phone
  - type: sms
stage: # 流程步骤
  - build
  - test
config: # 详细配置信息
  - stage: build
    name: template 模板测试
    template: 
      src: template.service.j2 
      dest: /tmp/docker.service
  - stage: test
    name: 上传文件
    copy: 
      src: "{{ .item }}"
      dest: /tmp/{{ .item }}
    with_items:
      - LICENSE
    tags: # 指定主机执行
      - 192.168.0.10
  - stage: what
    name: 非stage测试
    shell: whoami
  - stage: build
    name: 获取go version
    shell: lsb_release -a
  - stage: test
    name: 获取主机名
    shell: "{{.item}}"
    with_items:
    - hostname
    - ip a|grep eth0
    - pwd
    - uname -a
    - docker ps && 
      docker images
    tags:
      - 192.168.0.250
  - stage: test
    name: 查看docker信息
    shell: systemctl status sshd
`

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "查看完整yaml配置",
	Long: `yaml配置模板：
1. host 目标主机
2. remote_user 远程用户
3. roles 执行具体stage
4. vars 全局环境变量
5. hooks 消息发送
6. stage 定义流程步骤
7. config 详细配置信息`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(yaml)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
