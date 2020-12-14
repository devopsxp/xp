# xp
xp is my name,but this project is look like ansible-playbook and pipeline together

# 介绍

该项目主要想实现流水线+自动化实现工作和devops都能适应，两者互补。采用技术栈借鉴：

- [x] Ansbile Playbook
- [x] Gitlab CI

插件接口流程：

* pipeline
    * User 用户管理
    * Host 主机清单
    * Playbook
        * YAML
        * module
    * Plugin
        * start
        * stop
        * status
        * init
    * input
        * host conn check
        * host env
        * yaml module 分析
    * filter
        * 执行各个module
            * Role Module
                * copy
                * shell
                * template
                * ssh
                * docker
                * k8s
                * 网络设备
                * snmp等
                * 测试
                  * k6 run test.js
            * 执行
                * RPC
                * RESTFULL
                * CLI
        * 管理执行的生命周期
            * prepare
            * before
            * runtime
            * after
    * output
        * 输出结果
        * 返回状态

# 功能模块

- [x] yaml解析(cobra viper支持)
  - [ ] 获取环境变量，包括自定义、目标主机基本信息、本地ENV信息等
  - [x] 识别shell_role中shell、copy、template等三级模块
  - [x] go template support
  - [x] docker support
  - [x] systemd support（服务管理）
  - [x] k8s support
    - [x] [client-go](github.com/lflxp/lflxp-kubectl)
    - [x] init containerd
      - [x] git clone
    - [x] WorkingDir
    - [x] Shell+Command+Args
    - [ ] Env
    - [ ] Volume
      - [ ] HostPath
      - [x] Empty
      - [ ] PV/PVC
      - [ ] StorageClass
    - [ ] containerd
      - [ ] lifecycle
        - [ ] postStart
        - [ ] type: Sidecar
      - [x] 顺序执行
    - [x] 自动Delete Complete Pod
      - [x] pipeline hooks+output plugin调用k8s api进行删除
  - [x] ssh + docker support
  - [x] Hooks钩子函数支撑output plugin输出，目前支持console、email、wechat，todo： phone、sms、elasticsearch、log
  - [x] tags 指定主机执行
  - [x] 动态环境变量(cobra支持)
  - [x] with_items迭代器
  - [x] 缓存中间产物
  - [x] include命令，允许导入复杂yaml文件夹的大量引用，类似ansible-playbook roles
  - [x] fetch模块
  - [x] user模块
  - [x] git模块
  - [ ] group模块
  - [x] copy模块
  - [ ] yum模块
  - [ ] file模块
  - [ ] setup模块
- [x] Debug日志
- [x] CLI命令行工具(cobra)
- [ ] 功能文件夹，提供：files、hosts、env等特殊目录模块
- [ ] roles ansible模块
  - [x] yaml目录通过include模块引入
- [ ] module man模块说明文档
- [x] module plugin插件机制
- [ ] ssh [连接功能](https://github.com/mojocn/felix)
- [x] 各个步骤的计时器和总执行计时
- [x] Retry重试机制
  - [x] 超时重试
  - [x] 错误重试
- [x] 消息发送
  - [x] 邮件
  - [x] 企业微信/叮叮 
  - [x] 短信
- [ ] 中间件对接
  - [ ] sonarque
  - [ ] jmeter
  - [ ] jenkins
  - [ ] 安全扫描
  - [ ] ArgoCD
- [ ] 改造计划
  - [ ] 微服务改造
    - [ ] RPC AGENT
    - [ ] RPC Server
  - [ ] CRD改造
    - [ ] 声明式任务流水线执行
  - [ ] 日志Call Back机制
    - [ ] back to server
    - [ ] back to es 
    - [ ] back to kafka
    - [ ] back to logstash
    - [ ] back to Fluentd/Filebeat
- [ ] 前端页面
  - [ ] 多租户
  - [ ] pipeline yaml管理
  - [ ] pipeline执行历史管理
  - [ ] 权限管理
  - [ ] CMDB
  - [ ] CI/CD管理
  - [ ] pipeline 可视化
  - [ ] git源代码管理
  - [ ] devops工具链对接
- [ ] 容器化
  - [ ] docker support
    - [x] yaml新增images字段
    - [ ] [Remote API](https://docs.docker.com/engine/api/v1.24/)
  - [ ] k8s support
    - [x] yaml新增k8s字段
    - [ ] k8s agent/operator
    - [ ] pod all in one
  - [ ] ssh + docker support
  - [ ] 中间产物缓存


# Useage

> go build
> ./xp test --config devopsxp.yaml

# Test

> make

## 测试执行流程

cli -> main.go -> root.go -> test.go -> pipeline -> init -> start -> check(ssh) -> input(localyaml) -> filter(shell) -> output(console) -> stop

# Module

`Remove Check Plugin`

## Input Plugin

- [x] localyaml

## Filter Plugin

- [x] shell
- [x] ssh

## Output Plugin

- [x] console

# 配置信息

本工具采用ssh免密登录进行远程主机命令的执行，需要ssh私钥进行连接，默对认获取文件地址为：~/.ssh/id_rsa

## Like Ansible Playbook YAML

### 目标主机配置

配置目标执行主机，支持ip端扫描执行。

```yaml
host: # 目标主机
  - 127.0.0.1
  - 192.168.0.1-20
```

### 远程执行用户

```yaml
remote_user: root
```

### pipeline管理与stage管理

* roles用于限制stage是否执行，后期可用根据roles来实现流程动态选取
* stage用于限制config中执行的顺序和stage，没有列入stage的一律不执行

```yaml
roles: # 执行具体stage
  - build
  - test
stage: # 流程步骤
  - build
  - test
```

### 内置环境变量

用于动态嵌入output输出内容，配合go template实现数据插入。

`todo`：后期考虑加入`本机环境变量`和`远程主机环境变量`，供流程提取数据进行执行操作。

`格式`: 支持任意格式的数据（符合yaml语法），解析采用go template进行，`注意`：变量使用{{.Status}}

```yaml
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
```

### Hooks 钩子函数参数配置

消息发送钩子函数参数,支持同时发送多端配置发送，支持类型：

- [x] console
- [x] email
- [x] wechat
- [ ] phone
- [ ] sms
- [x] log

```yaml
hooks: # 消息发送 全局一个 可以一次发送所有结果到多个渠道
  - type: console # 类型，支持：console|email|wechat|phone|sms
  - type: email # email邮箱类型
    alias: xp战队 # 邮箱昵称
    email_user: "xp@xp.com" # 发送邮箱帐号
    email_pwd: "xppwdxp" # 发送邮箱密码
    email_smtp: "smtp.exmail.qq.com" # 邮箱发送smtp服务器
    email_smtp_port: "465" # smtp服务器端口
    email_to: # 接收邮件人员
      - xptest@xp.com
    template: # 告警模板
      title: "告警title" # 标题
      text: "{{.Status}} {{.title}}
      {{range $key,$value := .logs}}
        {{range $k,$v := $value}}
        ------ {{$k}} {{$v}}<br/>
        {{end}}
      {{end}}" # 文本模板
      path: template.service.j2 # 模板文件路径（与文本模板二选一，同时存在优先选择文本模板）
      vars: # 内置固定参数
        Status: true
        title: "模板测试"
        serviceName: "xp"
  - type: wechat # wechat类型
    address: # 企业微信机器人地址 支持多机器人批量发送
      - https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=a-b-c-d-e
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
{{end}}' # 文本模板
      path: template.service.j2 # 模板文件路径（与文本模板二选一，同时存在优先选择文本模板）
      vars: # 内置固定参数（与全局vars共享数据，这里设置相同参数会覆盖全局key）
        title: "模板测试"
        serviceName: "xp"
  - type: phone
  - type: sms
```

### 详细stage config配置

这里是主要的逻辑编写单元，实现各种pipeline的编写和stage的区分，目前支持的模块有：

- shell 模块，执行shell命令
- template 模块，执行基于go template的模板渲染和远程上传
- copy 模块，执行本地文件上传到目标主机指定地址

```yaml
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
      - 2
      - 3
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
  - stage: test
    name: 查看docker信息
    shell: systemctl status sshd
```

## Template 模板示例，采用go template，value来自上面`YAML` template模块的vars选项

```go
{{if .Status}} <font color="info">[成功] </font>{{ .title}} {{- else}}<font color="warning">[失败] </font> {{ .title}}  {{- end}} 
> 服务名称：<font color="comment">{{.serviceName}}</font> 
流水线名称：<font color="comment">{{.pipelineName}}</font> 
流水线ID：<font color="comment">{{.pipelineId}}</font> 
流水线地址：<font color="comment">{{.pipelineUrl}}</font> 
构建分支：<font color="comment">{{.branch}}</font> 
阶段：<font color="comment">{{.stage}}</font> 
构建人：<font color="comment">{{.user}}</font> 
COMMITID：<font color="comment">{{.tags}}</font> 
COMMITINFO: <font color="comment">{{.info}}</font> 

{{range $e := .alerts}}
## {{if eq $e.status "firing"}} [<font color="warning">{{$e.labels.severity}}</font>] {{ $e.labels.node}} [监控详情]({{$e.generatorURL}})

>> 告警阈值: <font color="warning">{{$e.labels.threshold_value}}</font> 
当前数值: <font color="warning">{{$e.annotations.current_value}}</font> 
开始时间: <font color="comment">{{$e.startsAt}}</font> 
结束时间: <font color="comment">{{$e.endsAt}}</font> 

{{- else}}# [<font color="info">{{$e.labels.severity}}</font>] {{ $e.labels.node}} [监控详情]({{$e.generatorURL}})

>> 告警阈值: <font color="info">{{$e.labels.threshold_value}}</font> 
当前数值: <font color="info">{{$e.annotations.current_value}}</font> 
开始时间: <font color="comment">{{$e.startsAt}}</font> 
结束时间: <font color="comment">{{$e.endsAt}}</font> 

{{- end}}

{{end}}
```

# TODO 

1. cli+module+shell
  * inputPlugin
  * 匹配数据
  * pipeline
2. systemd 服务管理模块
  * 匹配目标主机os系统
  * 根据目标主机服务管理方式进行service管理服务启停
